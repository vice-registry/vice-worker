// Autogenerated by Thrift Compiler (0.10.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
        "flag"
        "fmt"
        "math"
        "net"
        "net/url"
        "os"
        "strconv"
        "strings"
        "git.apache.org/thrift.git/lib/go/thrift"
        "bwlp"
)


func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  int getVersion(int clientVersion)")
  fmt.Fprintln(os.Stderr, "  string getSupportedFeatures()")
  fmt.Fprintln(os.Stderr, "  SatelliteConfig getConfiguration()")
  fmt.Fprintln(os.Stderr, "  TransferInformation requestImageVersionUpload(Token userToken, UUID imageBaseId, i64 fileSize,  blockHashes, string machineDescription)")
  fmt.Fprintln(os.Stderr, "  void updateBlockHashes(Token uploadToken,  blockHashes)")
  fmt.Fprintln(os.Stderr, "  void cancelUpload(Token uploadToken)")
  fmt.Fprintln(os.Stderr, "  TransferStatus queryUploadStatus(Token uploadToken)")
  fmt.Fprintln(os.Stderr, "  TransferInformation requestDownload(Token userToken, UUID imageVersionId)")
  fmt.Fprintln(os.Stderr, "  void cancelDownload(string downloadToken)")
  fmt.Fprintln(os.Stderr, "  void isAuthenticated(Token userToken)")
  fmt.Fprintln(os.Stderr, "  WhoamiInfo whoami(Token userToken)")
  fmt.Fprintln(os.Stderr, "  void invalidateSession(Token userToken)")
  fmt.Fprintln(os.Stderr, "   getUserList(Token userToken, i32 page)")
  fmt.Fprintln(os.Stderr, "  SatelliteUserConfig getUserConfig(Token userToken)")
  fmt.Fprintln(os.Stderr, "  void setUserConfig(Token userToken, SatelliteUserConfig config)")
  fmt.Fprintln(os.Stderr, "   getOperatingSystems()")
  fmt.Fprintln(os.Stderr, "   getVirtualizers()")
  fmt.Fprintln(os.Stderr, "   getAllOrganizations()")
  fmt.Fprintln(os.Stderr, "   getLocations()")
  fmt.Fprintln(os.Stderr, "  SatelliteStatus getStatus()")
  fmt.Fprintln(os.Stderr, "   getImageList(Token userToken,  tagSearch, i32 page)")
  fmt.Fprintln(os.Stderr, "  ImageDetailsRead getImageDetails(Token userToken, UUID imageBaseId)")
  fmt.Fprintln(os.Stderr, "  UUID createImage(Token userToken, string imageName)")
  fmt.Fprintln(os.Stderr, "  void updateImageBase(Token userToken, UUID imageBaseId, ImageBaseWrite image)")
  fmt.Fprintln(os.Stderr, "  void updateImageVersion(Token userToken, UUID imageVersionId, ImageVersionWrite image)")
  fmt.Fprintln(os.Stderr, "  void deleteImageVersion(Token userToken, UUID imageVersionId)")
  fmt.Fprintln(os.Stderr, "  void deleteImageBase(Token userToken, UUID imageBaseId)")
  fmt.Fprintln(os.Stderr, "  void writeImagePermissions(Token userToken, UUID imageBaseId,  permissions)")
  fmt.Fprintln(os.Stderr, "   getImagePermissions(Token userToken, UUID imageBaseId)")
  fmt.Fprintln(os.Stderr, "  void setImageOwner(Token userToken, UUID imageBaseId, UUID newOwnerId)")
  fmt.Fprintln(os.Stderr, "  void setImageVersionExpiry(Token userToken, UUID imageBaseId, UnixTimestamp expireTime)")
  fmt.Fprintln(os.Stderr, "  string getImageVersionVirtConfig(Token userToken, UUID imageVersionId)")
  fmt.Fprintln(os.Stderr, "  void setImageVersionVirtConfig(Token userToken, UUID imageVersionId, string meta)")
  fmt.Fprintln(os.Stderr, "  UUID requestImageReplication(Token userToken, UUID imageVersionId)")
  fmt.Fprintln(os.Stderr, "  UUID publishImageVersion(Token userToken, UUID imageVersionId)")
  fmt.Fprintln(os.Stderr, "  UUID createLecture(Token userToken, LectureWrite lecture)")
  fmt.Fprintln(os.Stderr, "  void updateLecture(Token userToken, UUID lectureId, LectureWrite lecture)")
  fmt.Fprintln(os.Stderr, "   getLectureList(Token userToken, i32 page)")
  fmt.Fprintln(os.Stderr, "  LectureRead getLectureDetails(Token userToken, UUID lectureId)")
  fmt.Fprintln(os.Stderr, "  void deleteLecture(Token userToken, UUID lectureId)")
  fmt.Fprintln(os.Stderr, "  void writeLecturePermissions(Token userToken, UUID lectureId,  permissions)")
  fmt.Fprintln(os.Stderr, "   getLecturePermissions(Token userToken, UUID lectureId)")
  fmt.Fprintln(os.Stderr, "  void setLectureOwner(Token userToken, UUID lectureId, UUID newOwnerId)")
  fmt.Fprintln(os.Stderr)
  os.Exit(0)
}

func main() {
  flag.Usage = Usage
  var host string
  var port int
  var protocol string
  var urlString string
  var framed bool
  var useHttp bool
  var parsedUrl url.URL
  var trans thrift.TTransport
  _ = strconv.Atoi
  _ = math.Abs
  flag.Usage = Usage
  flag.StringVar(&host, "h", "localhost", "Specify host and port")
  flag.IntVar(&port, "p", 9090, "Specify port")
  flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
  flag.StringVar(&urlString, "u", "", "Specify the url")
  flag.BoolVar(&framed, "framed", false, "Use framed transport")
  flag.BoolVar(&useHttp, "http", false, "Use http")
  flag.Parse()
  
  if len(urlString) > 0 {
    parsedUrl, err := url.Parse(urlString)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
    host = parsedUrl.Host
    useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
  } else if useHttp {
    _, err := url.Parse(fmt.Sprint("http://", host, ":", port))
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
  }
  
  cmd := flag.Arg(0)
  var err error
  if useHttp {
    trans, err = thrift.NewTHttpClient(parsedUrl.String())
  } else {
    portStr := fmt.Sprint(port)
    if strings.Contains(host, ":") {
           host, portStr, err = net.SplitHostPort(host)
           if err != nil {
                   fmt.Fprintln(os.Stderr, "error with host:", err)
                   os.Exit(1)
           }
    }
    trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
    if err != nil {
      fmt.Fprintln(os.Stderr, "error resolving address:", err)
      os.Exit(1)
    }
    if framed {
      trans = thrift.NewTFramedTransport(trans)
    }
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating transport", err)
    os.Exit(1)
  }
  defer trans.Close()
  var protocolFactory thrift.TProtocolFactory
  switch protocol {
  case "compact":
    protocolFactory = thrift.NewTCompactProtocolFactory()
    break
  case "simplejson":
    protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
    break
  case "json":
    protocolFactory = thrift.NewTJSONProtocolFactory()
    break
  case "binary", "":
    protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
    Usage()
    os.Exit(1)
  }
  client := bwlp.NewSatelliteServerClientFactory(trans, protocolFactory)
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "getVersion":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "GetVersion requires 1 args")
      flag.Usage()
    }
    argvalue0, err128 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err128 != nil {
      Usage()
      return
    }
    value0 := bwlp.Int(argvalue0)
    fmt.Print(client.GetVersion(value0))
    fmt.Print("\n")
    break
  case "getSupportedFeatures":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetSupportedFeatures requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetSupportedFeatures())
    fmt.Print("\n")
    break
  case "getConfiguration":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetConfiguration requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetConfiguration())
    fmt.Print("\n")
    break
  case "requestImageVersionUpload":
    if flag.NArg() - 1 != 5 {
      fmt.Fprintln(os.Stderr, "RequestImageVersionUpload requires 5 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    argvalue2, err131 := (strconv.ParseInt(flag.Arg(3), 10, 64))
    if err131 != nil {
      Usage()
      return
    }
    value2 := argvalue2
    arg132 := flag.Arg(4)
    mbTrans133 := thrift.NewTMemoryBufferLen(len(arg132))
    defer mbTrans133.Close()
    _, err134 := mbTrans133.WriteString(arg132)
    if err134 != nil { 
      Usage()
      return
    }
    factory135 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt136 := factory135.GetProtocol(mbTrans133)
    containerStruct3 := bwlp.NewSatelliteServerRequestImageVersionUploadArgs()
    err137 := containerStruct3.ReadField4(jsProt136)
    if err137 != nil {
      Usage()
      return
    }
    argvalue3 := containerStruct3.BlockHashes
    value3 := argvalue3
    argvalue4 := []byte(flag.Arg(5))
    value4 := argvalue4
    fmt.Print(client.RequestImageVersionUpload(value0, value1, value2, value3, value4))
    fmt.Print("\n")
    break
  case "updateBlockHashes":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "UpdateBlockHashes requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    arg140 := flag.Arg(2)
    mbTrans141 := thrift.NewTMemoryBufferLen(len(arg140))
    defer mbTrans141.Close()
    _, err142 := mbTrans141.WriteString(arg140)
    if err142 != nil { 
      Usage()
      return
    }
    factory143 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt144 := factory143.GetProtocol(mbTrans141)
    containerStruct1 := bwlp.NewSatelliteServerUpdateBlockHashesArgs()
    err145 := containerStruct1.ReadField2(jsProt144)
    if err145 != nil {
      Usage()
      return
    }
    argvalue1 := containerStruct1.BlockHashes
    value1 := argvalue1
    fmt.Print(client.UpdateBlockHashes(value0, value1))
    fmt.Print("\n")
    break
  case "cancelUpload":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "CancelUpload requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    fmt.Print(client.CancelUpload(value0))
    fmt.Print("\n")
    break
  case "queryUploadStatus":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "QueryUploadStatus requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    fmt.Print(client.QueryUploadStatus(value0))
    fmt.Print("\n")
    break
  case "requestDownload":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "RequestDownload requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.RequestDownload(value0, value1))
    fmt.Print("\n")
    break
  case "cancelDownload":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "CancelDownload requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.CancelDownload(value0))
    fmt.Print("\n")
    break
  case "isAuthenticated":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "IsAuthenticated requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    fmt.Print(client.IsAuthenticated(value0))
    fmt.Print("\n")
    break
  case "whoami":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "Whoami requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    fmt.Print(client.Whoami(value0))
    fmt.Print("\n")
    break
  case "invalidateSession":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "InvalidateSession requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    fmt.Print(client.InvalidateSession(value0))
    fmt.Print("\n")
    break
  case "getUserList":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "GetUserList requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    tmp1, err155 := (strconv.Atoi(flag.Arg(2)))
    if err155 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    fmt.Print(client.GetUserList(value0, value1))
    fmt.Print("\n")
    break
  case "getUserConfig":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "GetUserConfig requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    fmt.Print(client.GetUserConfig(value0))
    fmt.Print("\n")
    break
  case "setUserConfig":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "SetUserConfig requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    arg158 := flag.Arg(2)
    mbTrans159 := thrift.NewTMemoryBufferLen(len(arg158))
    defer mbTrans159.Close()
    _, err160 := mbTrans159.WriteString(arg158)
    if err160 != nil {
      Usage()
      return
    }
    factory161 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt162 := factory161.GetProtocol(mbTrans159)
    argvalue1 := bwlp.NewSatelliteUserConfig()
    err163 := argvalue1.Read(jsProt162)
    if err163 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    fmt.Print(client.SetUserConfig(value0, value1))
    fmt.Print("\n")
    break
  case "getOperatingSystems":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetOperatingSystems requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetOperatingSystems())
    fmt.Print("\n")
    break
  case "getVirtualizers":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetVirtualizers requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetVirtualizers())
    fmt.Print("\n")
    break
  case "getAllOrganizations":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetAllOrganizations requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetAllOrganizations())
    fmt.Print("\n")
    break
  case "getLocations":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetLocations requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetLocations())
    fmt.Print("\n")
    break
  case "getStatus":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetStatus requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetStatus())
    fmt.Print("\n")
    break
  case "getImageList":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "GetImageList requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    arg165 := flag.Arg(2)
    mbTrans166 := thrift.NewTMemoryBufferLen(len(arg165))
    defer mbTrans166.Close()
    _, err167 := mbTrans166.WriteString(arg165)
    if err167 != nil { 
      Usage()
      return
    }
    factory168 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt169 := factory168.GetProtocol(mbTrans166)
    containerStruct1 := bwlp.NewSatelliteServerGetImageListArgs()
    err170 := containerStruct1.ReadField2(jsProt169)
    if err170 != nil {
      Usage()
      return
    }
    argvalue1 := containerStruct1.TagSearch
    value1 := argvalue1
    tmp2, err171 := (strconv.Atoi(flag.Arg(3)))
    if err171 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    fmt.Print(client.GetImageList(value0, value1, value2))
    fmt.Print("\n")
    break
  case "getImageDetails":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "GetImageDetails requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.GetImageDetails(value0, value1))
    fmt.Print("\n")
    break
  case "createImage":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "CreateImage requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := argvalue1
    fmt.Print(client.CreateImage(value0, value1))
    fmt.Print("\n")
    break
  case "updateImageBase":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "UpdateImageBase requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    arg178 := flag.Arg(3)
    mbTrans179 := thrift.NewTMemoryBufferLen(len(arg178))
    defer mbTrans179.Close()
    _, err180 := mbTrans179.WriteString(arg178)
    if err180 != nil {
      Usage()
      return
    }
    factory181 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt182 := factory181.GetProtocol(mbTrans179)
    argvalue2 := bwlp.NewImageBaseWrite()
    err183 := argvalue2.Read(jsProt182)
    if err183 != nil {
      Usage()
      return
    }
    value2 := argvalue2
    fmt.Print(client.UpdateImageBase(value0, value1, value2))
    fmt.Print("\n")
    break
  case "updateImageVersion":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "UpdateImageVersion requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    arg186 := flag.Arg(3)
    mbTrans187 := thrift.NewTMemoryBufferLen(len(arg186))
    defer mbTrans187.Close()
    _, err188 := mbTrans187.WriteString(arg186)
    if err188 != nil {
      Usage()
      return
    }
    factory189 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt190 := factory189.GetProtocol(mbTrans187)
    argvalue2 := bwlp.NewImageVersionWrite()
    err191 := argvalue2.Read(jsProt190)
    if err191 != nil {
      Usage()
      return
    }
    value2 := argvalue2
    fmt.Print(client.UpdateImageVersion(value0, value1, value2))
    fmt.Print("\n")
    break
  case "deleteImageVersion":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "DeleteImageVersion requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.DeleteImageVersion(value0, value1))
    fmt.Print("\n")
    break
  case "deleteImageBase":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "DeleteImageBase requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.DeleteImageBase(value0, value1))
    fmt.Print("\n")
    break
  case "writeImagePermissions":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "WriteImagePermissions requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    arg198 := flag.Arg(3)
    mbTrans199 := thrift.NewTMemoryBufferLen(len(arg198))
    defer mbTrans199.Close()
    _, err200 := mbTrans199.WriteString(arg198)
    if err200 != nil { 
      Usage()
      return
    }
    factory201 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt202 := factory201.GetProtocol(mbTrans199)
    containerStruct2 := bwlp.NewSatelliteServerWriteImagePermissionsArgs()
    err203 := containerStruct2.ReadField3(jsProt202)
    if err203 != nil {
      Usage()
      return
    }
    argvalue2 := containerStruct2.Permissions
    value2 := argvalue2
    fmt.Print(client.WriteImagePermissions(value0, value1, value2))
    fmt.Print("\n")
    break
  case "getImagePermissions":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "GetImagePermissions requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.GetImagePermissions(value0, value1))
    fmt.Print("\n")
    break
  case "setImageOwner":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "SetImageOwner requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    argvalue2 := flag.Arg(3)
    value2 := bwlp.UUID(argvalue2)
    fmt.Print(client.SetImageOwner(value0, value1, value2))
    fmt.Print("\n")
    break
  case "setImageVersionExpiry":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "SetImageVersionExpiry requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    argvalue2, err211 := (strconv.ParseInt(flag.Arg(3), 10, 64))
    if err211 != nil {
      Usage()
      return
    }
    value2 := bwlp.UnixTimestamp(argvalue2)
    fmt.Print(client.SetImageVersionExpiry(value0, value1, value2))
    fmt.Print("\n")
    break
  case "getImageVersionVirtConfig":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "GetImageVersionVirtConfig requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.GetImageVersionVirtConfig(value0, value1))
    fmt.Print("\n")
    break
  case "setImageVersionVirtConfig":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "SetImageVersionVirtConfig requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    argvalue2 := []byte(flag.Arg(3))
    value2 := argvalue2
    fmt.Print(client.SetImageVersionVirtConfig(value0, value1, value2))
    fmt.Print("\n")
    break
  case "requestImageReplication":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "RequestImageReplication requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.RequestImageReplication(value0, value1))
    fmt.Print("\n")
    break
  case "publishImageVersion":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "PublishImageVersion requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.PublishImageVersion(value0, value1))
    fmt.Print("\n")
    break
  case "createLecture":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "CreateLecture requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    arg222 := flag.Arg(2)
    mbTrans223 := thrift.NewTMemoryBufferLen(len(arg222))
    defer mbTrans223.Close()
    _, err224 := mbTrans223.WriteString(arg222)
    if err224 != nil {
      Usage()
      return
    }
    factory225 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt226 := factory225.GetProtocol(mbTrans223)
    argvalue1 := bwlp.NewLectureWrite()
    err227 := argvalue1.Read(jsProt226)
    if err227 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    fmt.Print(client.CreateLecture(value0, value1))
    fmt.Print("\n")
    break
  case "updateLecture":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "UpdateLecture requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    arg230 := flag.Arg(3)
    mbTrans231 := thrift.NewTMemoryBufferLen(len(arg230))
    defer mbTrans231.Close()
    _, err232 := mbTrans231.WriteString(arg230)
    if err232 != nil {
      Usage()
      return
    }
    factory233 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt234 := factory233.GetProtocol(mbTrans231)
    argvalue2 := bwlp.NewLectureWrite()
    err235 := argvalue2.Read(jsProt234)
    if err235 != nil {
      Usage()
      return
    }
    value2 := argvalue2
    fmt.Print(client.UpdateLecture(value0, value1, value2))
    fmt.Print("\n")
    break
  case "getLectureList":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "GetLectureList requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    tmp1, err237 := (strconv.Atoi(flag.Arg(2)))
    if err237 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    fmt.Print(client.GetLectureList(value0, value1))
    fmt.Print("\n")
    break
  case "getLectureDetails":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "GetLectureDetails requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.GetLectureDetails(value0, value1))
    fmt.Print("\n")
    break
  case "deleteLecture":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "DeleteLecture requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.DeleteLecture(value0, value1))
    fmt.Print("\n")
    break
  case "writeLecturePermissions":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "WriteLecturePermissions requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    arg244 := flag.Arg(3)
    mbTrans245 := thrift.NewTMemoryBufferLen(len(arg244))
    defer mbTrans245.Close()
    _, err246 := mbTrans245.WriteString(arg244)
    if err246 != nil { 
      Usage()
      return
    }
    factory247 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt248 := factory247.GetProtocol(mbTrans245)
    containerStruct2 := bwlp.NewSatelliteServerWriteLecturePermissionsArgs()
    err249 := containerStruct2.ReadField3(jsProt248)
    if err249 != nil {
      Usage()
      return
    }
    argvalue2 := containerStruct2.Permissions
    value2 := argvalue2
    fmt.Print(client.WriteLecturePermissions(value0, value1, value2))
    fmt.Print("\n")
    break
  case "getLecturePermissions":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "GetLecturePermissions requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    fmt.Print(client.GetLecturePermissions(value0, value1))
    fmt.Print("\n")
    break
  case "setLectureOwner":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "SetLectureOwner requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := bwlp.Token(argvalue0)
    argvalue1 := flag.Arg(2)
    value1 := bwlp.UUID(argvalue1)
    argvalue2 := flag.Arg(3)
    value2 := bwlp.UUID(argvalue2)
    fmt.Print(client.SetLectureOwner(value0, value1, value2))
    fmt.Print("\n")
    break
  case "":
    Usage()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
  }
}
