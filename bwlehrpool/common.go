package bwlehrpool

import (
	"log"
	"github.com/vice-registry/vice-api/models"
	"github.com/OpenSLX/bwlp-go-client/client"
)

func Login(environment *models.Environment) (*client.SessionHandler, error) {
	//endpoint := environment.Credentials.Location
	username := environment.Credentials.Username
	password := environment.Credentials.Password
	// TODO extract the data from environment.Credentials.Location
	// TODO use Credentials.Specifics for masterserver ports 
	endpointData := client.MasterServerEndpoint{
					Hostname: "bwlp-masterserver.ruf.uni-freiburg.de",
					PortSSL: 9091,
					PortPlain: 9090,
	}
	// get a new handler first
	handler, err := client.NewSessionHandler(&endpointData)
	if err != nil {
		log.Printf("Error initializing bwLehrpool handler: %s", err)
		return nil, err
	}
	// perform login
	if err := handler.Login(username, password); err != nil {
		log.Printf("Error logging in with bwLehrpool: %s", err)
		return nil, err
	}
	return handler, nil
}
