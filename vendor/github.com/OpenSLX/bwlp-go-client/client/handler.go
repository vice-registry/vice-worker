package client

import (
	"fmt"
	"log"
	"sync"
	"errors"

	"github.com/OpenSLX/bwlp-go-client/bwlp"
)
var (
	singleMaster sync.Once
	singleSatellite sync.Once
)
type SessionHandler struct {
	// master endpoint + client
	masterEndpoint *ServerEndpoint
	masterClient *bwlp.MasterServerClient
	// satellite endpoint + client
	satEndpoint *ServerEndpoint
	satClient *bwlp.SatelliteServerClient
	SessionData *bwlp.ClientSessionData
}

func NewSessionHandler(masterServerEndpoint *ServerEndpoint) (*SessionHandler) {
	if masterServerEndpoint == nil {
		log.Printf("No endpoint to a masterserver was given!\n")
		return nil
	}
	newHandler := SessionHandler{
		masterEndpoint: masterServerEndpoint,
		masterClient: nil,
		satEndpoint: nil,
		satClient: nil,
		SessionData: nil,
	}
	return &newHandler
}

// Global access to the singleton masterserver client instance
func (handler *SessionHandler) GetMasterClient() (*bwlp.MasterServerClient) {
	// initialize the client only once, in essence
	// a simple kind of singleton pattern
	singleMaster.Do(func() {
		masterServerAddress := fmt.Sprintf("%s:%d", handler.masterEndpoint.Hostname, handler.masterEndpoint.PortSSL)
		client, err := initMasterClient(masterServerAddress)
		if err != nil {
			log.Printf("Error initialising master client: %s\n", err)
		}
		handler.masterClient = client
	})
	if handler.masterClient == nil {
		// TODO handle dead clients
		log.Printf("Masterserver client lost connection!\n")
	}
	return handler.masterClient
}

// Helper to create a new masterclient (currently not used, will be needed when handling client's lost connections)
func (handler *SessionHandler) newMasterClient() (*bwlp.MasterServerClient) {
		masterServerAddress := fmt.Sprintf("%s:%d", handler.masterEndpoint.Hostname, handler.masterEndpoint.PortSSL)
		client, err := initMasterClient(masterServerAddress)
		if err != nil {
			log.Printf("Error initialising master client: %s\n", err)
			return nil
		}
		return client
}

// Global setter for the satellite endpoint as this is dictated by the masterserver
// (or potentially chosen by the user in the future)
func (handler *SessionHandler) SetSatEndpoint(param *ServerEndpoint) error {
	if handler.satEndpoint != nil {
		log.Printf("Satellite server endpoint is already set.\n")
		return nil
	}
	if param == nil {
		return errors.New("Invalid endpoint given!")
	}
	// TODO user-supplied endpoints should be validated
	handler.satEndpoint = param
	return nil
}

// Global access to the singleton satellite client instance
func (handler *SessionHandler) GetSatClient() (*bwlp.SatelliteServerClient) {
	// initialize the client only once, in essence
	// a simple kind of singleton pattern
	singleSatellite.Do(func() {
		if handler.satEndpoint == nil {
			log.Printf("Failed to retrieve client - satellite endpoint not set!")
			return
		}
		satServerAddress := fmt.Sprintf("%s:%d", handler.satEndpoint.Hostname, handler.satEndpoint.PortSSL)
		client, err := initSatClient(satServerAddress)
		if err != nil {
			log.Printf("Error initialising sat client: %s\n", err)
		}
		handler.satClient = client
	})
	if handler.satClient == nil {
		// TODO handle dead clients
		log.Printf("Satellite client lost connection?")
	}
	return handler.satClient
}
