package actions

import (
	"fmt"
	"log"
	"strings"

	"github.com/vice-registry/vice-api/models"
	"github.com/vice-registry/vice-api/persistence"
	"github.com/vice-registry/vice-worker/bwlehrpool"
	"github.com/vice-registry/vice-worker/common"
	"github.com/vice-registry/vice-worker/openstack"
)

// Action encapsulates necessary information of an action
type Action struct {
	reference   string
	image       *models.Image
	deployment  *models.Deployment
	environment *models.Environment
	adaptor     adaptors.Adaptor
}

func handleAction(action Action) error {

	// select adaptor for handling action
	managementSoftware, err := getManagementSoftware(&action)
	if err != nil {
		log.Printf("Error in import handleAction: %s", err)
		return err
	}
	var adaptor adaptors.Adaptor
	switch strings.ToLower(managementSoftware) {
	case "openstack":
		adaptor = openstack.OpenStack{}
	case "bwlehrpool":
		adaptor = bwlehrpool.BwLehrpool{}
	case "docker":
		// TODO
	default:
		err := fmt.Errorf("no adaptor found for management layer %s", managementSoftware)
		log.Printf("Error in handleAction: %s", err)
		return err
	}
	action.adaptor = adaptor

	// select WorkerType and trigger adaptor
	if adaptors.WorkerType == adaptors.WorkerTypeImport {
		handleImport(action)
	} else if adaptors.WorkerType == adaptors.WorkerTypeExport {
		handleExport(action)
	} else {
		err := fmt.Errorf("unknown workertype %s", adaptors.WorkerType)
		return err
	}

	return nil

}

func getManagementSoftware(action *Action) (string, error) {
	if adaptors.WorkerType == adaptors.WorkerTypeImport {
		image, err := persistence.GetImage(action.reference)
		if err != nil {
			return "", err
		}
		action.image = image
		action.environment = image.OriginEnvironment
		return image.OriginEnvironment.Managementlayer.Software, nil
	} else if adaptors.WorkerType == adaptors.WorkerTypeExport {
		deployment, err := persistence.GetDeployment(action.reference)
		if err != nil {
			return "", err
		}
		action.deployment = deployment
		environment, err := persistence.GetEnvironment(deployment.EnvironmentID)
		if err != nil {
			return "", err
		}
		action.environment = environment
		return environment.Managementlayer.Software, nil
	} else {
		err := fmt.Errorf("unknown workertype %s", adaptors.WorkerType)
		return "", err
	}
}
