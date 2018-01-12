package client

import (
	"log"
	"crypto/tls"
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/OpenSLX/bwlp-go-client/bwlp"
)

type ServerEndpoint struct {
	Hostname string
	PortSSL int
	PortPlain int
}

func createTLSTransport(addr string) (*thrift.TTransport) {
	var transport thrift.TTransport
	cfg := &tls.Config{
		MinVersion:	tls.VersionTLS12,
		CurvePreferences:	[]tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites:	[]uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		InsecureSkipVerify: true,
		PreferServerCipherSuites: true,
	}
	transport, err := thrift.NewTSSLSocket(addr, cfg)
	if err != nil {
		log.Printf("Error opening SSL socket: %s\n", err)
		return nil
	}
	// framed transport is required
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport = transportFactory.GetTransport(transport)
	if err := transport.Open(); err != nil {
		log.Printf("Error opening transport layer for reading/writing: %s\n", err)
		return nil
	}
	return &transport
}

// Initialize the masterserver client using the server's
// expected transport (framed) and protocol (binary).
// Enforces the use of SSL for now.
func initMasterClient(addr string) (*bwlp.MasterServerClient, error) {
	transport := createTLSTransport(addr)
	if transport == nil {
		return nil, errors.New("Could not create TLS transport.")
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	// now retrieve a new client and test it
	masterClient := bwlp.NewMasterServerClientFactory((*transport), protocolFactory)
	if masterClient == nil {
		return nil, errors.New("Thrift client factory return nil client!")
	}

	if _, err := masterClient.Ping(); err != nil {
		log.Printf("Error pinging masterserver: %s\n", err)
		return nil, err
  }
	log.Printf("## Connection established to master: %s ##\n", addr)
	return masterClient, nil
}

// Initialize the satellite client similarly to the masterserver client
func initSatClient(addr string) (*bwlp.SatelliteServerClient, error) {
	transport := createTLSTransport(addr)
	if transport == nil {
		return nil, errors.New("Could not create TLS transport.")
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	// now retrieve a new client and test it
	satClient := bwlp.NewSatelliteServerClientFactory((*transport), protocolFactory)
	if satClient == nil {
		return nil, errors.New("Thrift client factory return nil client!")
	}
	if _, err := satClient.GetSupportedFeatures(); err != nil {
		log.Printf("Error testing sat client: %s\n", err)
		return nil, err
  }
	log.Printf("## Connection established to satellite: %s ##\n", addr)
	return satClient, nil
}
