package client

import (
	"log"
	"net"
	"io"
	"bufio"
	"strconv"
	"strings"
	"encoding/binary"
	"fmt"
	"github.com/OpenSLX/bwlp-go-client/bwlp"
)

type Transfer struct {
	transferTo bool
	data []byte
	conn net.Conn
	connReader *bufio.Reader
	connWriter *bufio.Writer

	fileSize int64
	chunkSize int64
	startOffset int64
	endOffset int64
	totalTransferred int64

	Ti *bwlp.TransferInformation
}

func NewTransfer(uploading bool, hostname string, ti *bwlp.TransferInformation, fileSize int64) *Transfer {
	// initialize connection
	conn, err := net.Dial("tcp", hostname + ":" + strconv.FormatInt(int64(ti.PlainPort), 10))
	if err != nil {
		log.Printf("Error establishing connection: %s\n", err)
		return nil
	}
	// init reader and writer
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	if reader == nil || writer == nil {
		log.Printf("Error initializing reader/writer for transfer.\n")
		return nil
	}
	// transferTo holds whether this is an upload (true) or a download (false)
	action := []byte("D")
	if uploading {
		action = []byte("U")
	}
	err = binary.Write(writer, binary.BigEndian, action)
	if err != nil {
		log.Printf("Error sending action initiator: %s\n", err)
		return nil
	}
	// send token
	if err := sendKeyValue(writer, "TOKEN", string(ti.Token)); err != nil {
		log.Printf("Failed to send token: %s\n", err)
		return nil
	}
	var chunkSize int64 = 16 * 1024 * 1024 // 16MB
	if chunkSize > fileSize {
		chunkSize = fileSize
	}
	t := Transfer{
		transferTo: uploading,
		conn: conn,
		connReader: reader,
		connWriter: writer,
		fileSize: fileSize,
		chunkSize: chunkSize,
		startOffset: 0,
		endOffset: chunkSize,
		totalTransferred: 0,
		Ti: ti,
	}
	return &t
}

func sendEndOfMeta(writer *bufio.Writer) (error) {
	// End of meta consists of two null bytes
	if err := binary.Write(writer, binary.BigEndian, []byte{0x00, 0x00}); err != nil {
		log.Printf("Error writing terminating sequence!")
		return err
	}
	writer.Flush()
	return nil
}

func sendKeyValue(writer *bufio.Writer, key string, value string) (error) {
	if len(key) <= 0 {
		return fmt.Errorf("Empty key!")
	}
	msg := key + "=" + value
	// To comply to java's readUtf8 method, we need to prepend the message
	// with the byte-encoded length of the message
	if err := binary.Write(writer, binary.BigEndian, int16(len(msg))); err != nil {
		log.Printf("Failed to write length!")
		return err
	}
	if err := binary.Write(writer, binary.BigEndian, []byte(msg)); err != nil {
		log.Println("Failed to write payload!", err)
		return err
	}
	return sendEndOfMeta(writer)
}

func readMetaData(reader *bufio.Reader) (string, error) {
	// first 2 bytes in java's modified UTF-8 contain the length
	metaLengthAsBytes := make([]byte, 2)
	if err := binary.Read(reader, binary.BigEndian, &metaLengthAsBytes); err != nil {
		log.Printf("Failed to read meta message length: %s\n", err)
		return "", err
	}
	metaLength := binary.BigEndian.Uint16(metaLengthAsBytes)
	metaBytes := make([]byte, metaLength)
	if err := binary.Read(reader, binary.BigEndian, metaBytes); err != nil {
		log.Printf("Failed to read actual meta message: %s\n", err)
		return "", err
	}
	if err := readEndOfMeta(reader); err != nil {
		log.Printf("%s", err)
		return "", err
	}
	return string(metaBytes[:]), nil
}

func readEndOfMeta(reader *bufio.Reader) error {
	readBytes := make([]byte, 2)
	if err := binary.Read(reader, binary.BigEndian, readBytes); err != nil {
		log.Println("Error reading terminating sequence!", err)
		return err
	}
	if len(readBytes) != 2 || readBytes[0] != 0x00 || readBytes[1] != 0x00 {
		return fmt.Errorf("Terminating sequence expected, but got: [% x]", readBytes)
	}
	return nil
}

func (t *Transfer) processRangeRequest() error {
	// read incoming range request
	meta, err := readMetaData(t.connReader)
	if err != nil {
		log.Printf("Error reading range: %s\n", err)
		return err
	}
	// expect RANGE=x:y and split it to get start = x and end = y
	metaKV := strings.Split(meta, "=")
	if metaKV[0] != "RANGE" {
		return fmt.Errorf("Invalid key received. Expected 'RANGE', got %s (%s)\n", metaKV[0], metaKV[1])
	}
	rangeBounds := strings.Split(metaKV[1], ":")
	start, err := strconv.ParseInt(rangeBounds[0], 10, 64)
	if err != nil {
		log.Printf("Error parsing starting range offset: %s\n", err)
		return err
	}
	end, err := strconv.ParseInt(rangeBounds[1], 10, 64)
	if err != nil {
		log.Printf("Error parsing ending range offset: %s\n", err)
		return err
	}
	// send confirmation
	err = sendKeyValue(t.connWriter, metaKV[0], metaKV[1])
	if err != nil {
		log.Printf("Error sending range confirmation: %s\n", err)
		return err
	}
	t.startOffset = start
	t.endOffset = end
	t.updateOffsets()
	return nil
}

func (t *Transfer) sendRangeRequest(startOffset int64, endOffset int64) error {
	// send range request
	rangeString := strconv.FormatInt(startOffset, 10) + ":" + strconv.FormatInt(endOffset, 10)
	if err := sendKeyValue(t.connWriter, "RANGE", rangeString); err != nil {
		return err
	}
	// read confirmation
	meta, err := readMetaData(t.connReader)
	if err != nil {
		log.Printf("Error reading range confirmation: %s\n", err)
		return err
	}
	// match?
	if meta != "RANGE=" + rangeString {
		return fmt.Errorf("Unexpected RANGE request response from server: %s", rangeString)
	}
	return nil
}
// Helper to update internal Transfer fields
func (t *Transfer) updateOffsets() {
	if t.totalTransferred == t.endOffset {
		t.startOffset = t.endOffset
		if t.endOffset + t.chunkSize > t.fileSize {
			t.endOffset = t.fileSize
		} else {
			t.endOffset = t.startOffset + t.chunkSize
		}
	}
}

// the read interface
func (t *Transfer) Read(p []byte) (n int, err error) {
	//log.Printf("TOTAL: %d / %d\n", t.totalTransferred, t.fileSize)
	if t.totalTransferred == t.fileSize {
		sendKeyValue(t.connWriter, "DONE", "")
		return 0, io.EOF
	}
	// check for buffered bytes from the remote connection
	if t.totalTransferred == 0 || t.totalTransferred == t.endOffset {
		t.updateOffsets()
		// no data cached, request from remote connection
		if err := t.sendRangeRequest(t.startOffset, t.endOffset); err != nil {
			return 0, err
		}
	}
	// now read data
	n, err = t.connReader.Read(p)
	if err != nil {
		if err != io.EOF {
			log.Printf("READ ERROR: %s\n", err)
		}
		return n, err
	}
	t.totalTransferred += int64(n)
	return n, nil
}

func (t *Transfer) Write(p []byte) (n int, err error) {
	// now we need to buffer bytes from the remote connection
	// before handing them out to the reader
	if t.totalTransferred == 0 || t.totalTransferred + int64(len(p)) > t.endOffset || t.totalTransferred == t.endOffset {
		// if the buffered writer has still buffered bytes, flush them
		// before checking for the next range message
		if t.connWriter.Buffered() != 0 {
			if err := t.connWriter.Flush(); err != nil {
				return 0, err
			}
		}
		// first range or current range was 
		if err := t.processRangeRequest(); err != nil {
			return 0, err
		}
	}
	// now write data
	n, err = t.connWriter.Write(p)
	if err != nil {
		if err != io.EOF {
			log.Printf("WRITE ERROR: %s\n", err)
		}
		return n, err
	}
	t.totalTransferred += int64(n)
	if t.totalTransferred == t.fileSize {
		if t.connWriter.Buffered() != 0 {
			if err := t.connWriter.Flush(); err != nil {
				return n, err
			}
		}

		// read done
		done, err := readMetaData(t.connReader)
		if err != nil {
			return n, err
		}
		if done == "DONE=" {
			return n, io.EOF
		}
		log.Printf("Received: %s\n", done)
	}
	return n, nil
}
