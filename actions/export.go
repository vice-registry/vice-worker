package actions

import (
	"log"

	"github.com/vice-registry/vice-api/persistence"
)

func handleExport(action Action) error {
	deployment := action.deployment
	log.Printf("handle export for deployment %+v", deployment)

	// call adaptor
	err := action.adaptor.Deploy(deployment)
	if err != nil {
		log.Printf("Error in handleExport for deploymentID %s: %s", action.reference, err)
		return err
	}

	// update deployment in couchbase
	deployment.EnvironmentReference = "" // TODO
	_, err = persistence.UpdateDeployment(deployment)
	if err != nil {
		log.Printf("Error in handleExport for deploymentID %s: %s", action.reference, err)
		return err
	}

	return nil
}
