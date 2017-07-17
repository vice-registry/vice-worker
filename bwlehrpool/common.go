package bwlehrpool

import (
	"fmt"
	"github.com/vice-registry/vice-api/models"
	"github.com/OpenSLX/bwlp-go-client/bwlp"
	"github.com/OpenSLX/bwlp-go-client/client"
)

type SessionHandler struct {
	Client *bwlp.MasterServerClient
	SessionData *bwlp.ClientSessionData
}

func NewSessionHandler() SessionHandler {
	return SessionHandler{Client: nil, SessionData: nil}
}

func (handler *SessionHandler) Login(environment *models.Environment) error {
	//endpoint := environment.Credentials.Location
	username := environment.Credentials.Username
	password := environment.Credentials.Password
	// TODO evaluate the need for Credentials.Specifics 
	// TODO extract the data from environment.Credentials.Location
	endpointData := client.MasterServerEndpoint{
					Hostname: "bwlp-masterserver.ruf.uni-freiburg.de",
					PortSSL: 9091,
					PortPlain: 9090,
	}
	// make sure we have a working client
	if err := handler.initClient(&endpointData); err != nil {
		fmt.Println("Error initializing masterserver client", err)
		return err
	}
	// perform login
	session, err := handler.Client.LocalAccountLogin(username, password)
	if err != nil {
		fmt.Printf("## Authentication failed: %s", err)
		return err
  }
	// store session data for later use
	handler.SessionData = session
	return nil
}

func (handler *SessionHandler) initClient(endpoint *client.MasterServerEndpoint) error {
	// already initialised?
	if handler.Client != nil {
		fmt.Println("Masterserver client already initialized.")
		return nil
	}
	// set environment's endpoint
	if err := client.SetEndpoint(endpoint); err != nil {
		fmt.Printf("Error setting endpoint during masterserver client initialisation: %s", err)
		return err
	}
	// get main masterserver client instance
	client := client.GetInstance()

	// verify that connection is established,
	_, err := client.Ping()
	if err != nil {
		fmt.Println("Error pinging masterserver:", err)
    return err
  }
	handler.Client = client
	return nil
}

// TEST functions 
func (handler *SessionHandler) GetPublicImages() error {
	iList, err := handler.Client.GetPublicImages(handler.SessionData.SessionId, 0)
	if err != nil {
		fmt.Println("Error fetching image list:", err)
		return err
	}
	for i := range iList {
		fmt.Println(iList[i])
	}
	return nil
}

func (handler *SessionHandler) GetImageDetails(imageBaseId *bwlp.UUID) error {
	iDetails, err := handler.Client.GetImageDetails(handler.SessionData.SessionId, *imageBaseId)
	if err != nil {
		fmt.Printf("Error fetching image details of '%s': %s", &imageBaseId, err)
		return err
	}
	fmt.Println(iDetails)
	return nil
}

