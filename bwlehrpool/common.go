package bwlehrpool

import (
	"log"
	"errors"

	"github.com/OpenSLX/bwlp-go-client/bwlp"
	"github.com/OpenSLX/bwlp-go-client/client"
	"github.com/vice-registry/vice-util/models"
)
type bwlpSpecifics struct {
	Details *bwlp.ImageDetailsRead
	Version *bwlp.ImageVersionDetails
	MachineDescription []byte
}

// The given environment's location must point to a bwlehrpool satellite server
func Login(environment *models.Environment) (*client.SessionHandler, error) {
	// Endpoint to central masterserver - not subject to change.
	masterEndpoint := client.ServerEndpoint{
		Hostname:  "bwlp-masterserver.ruf.uni-freiburg.de",
		PortSSL:   9091,
		PortPlain: 9090,
	}
	// get a new handler first
	handler := client.NewSessionHandler(&masterEndpoint)
	if handler == nil {
		return nil, errors.New("Error initializing bwLehrpool handler.")
	}
	// perform login with given credentials
	username := environment.Credentials.Username
	password := environment.Credentials.Password
	if err := handler.Login(username, password); err != nil {
		log.Printf("Error logging in with bwLehrpool: %s", err)
		return nil, err
	}

	// Now we set the satellite stored in the given environment
	satEndpoint := client.ServerEndpoint{
		Hostname:  environment.Credentials.Location,
		PortSSL:   9091,
		PortPlain: 9090,
	}
	// Note: this kinda goes against the original idea that the
	// masterserver gives us which satellites are registered for
	// this user - for demonstration purposes it will do however.
	if err := handler.SetSatEndpoint(&satEndpoint); err != nil {
		log.Printf("Failed to set satellite: %s\n", err)
		return nil, err
	}
	return handler, nil
}
