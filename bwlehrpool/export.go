package bwlehrpool

import (
	"log"
	"encoding/json"

	"github.com/vice-registry/vice-util/models"
	"github.com/vice-registry/vice-util/persistence"
	"github.com/vice-registry/vice-util/storeclient"
)

func handleExport(deployment *models.Deployment) error {
	log.Printf("Going to export deploymentID %s to bwLehrpool", deployment.ID)

	// get target environment to deploy to
	targetEnvironment, err := persistence.GetEnvironment(deployment.EnvironmentID)
	if err != nil {
		return err
	}

	// get source image to export
	sourceImage, err := persistence.GetImage(deployment.Imageid)
	if err != nil {
		return err
	}

	// get the image specifics to determine size (needed to by the uploader)
	var sourceImageSize int64 = -1
	var sourceImageMachineDescription []byte
	var sourceImageSpecifics bwlpSpecifics
	if err := json.Unmarshal([]byte(sourceImage.Specifics.(string)), &sourceImageSpecifics); err != nil {
		// if unmarshaling into bwlpSpecifics fails, the image came from another environment
		// TODO as a result, we need to determine the image size and machine description in some other way
		log.Printf("Failed to extract image specifics as bwlp image and version details.")
	} else {
		sourceImageSize = sourceImageSpecifics.Version.FileSize
		sourceImageMachineDescription = sourceImageSpecifics.MachineDescription
	}
	// now login with bwlehrpool's environment
	sessionHandler, err := Login(targetEnvironment)
	if err != nil {
		return err
	}

	// a new base image must be created before uploading a version to it
	destImageBaseId, err := sessionHandler.CreateImage(sourceImage.Title)
	if err != nil {
		return err
	}

	// now upload a version for the new base image
	uploader, err := sessionHandler.UploadImageVersion(destImageBaseId, sourceImageSize, sourceImageMachineDescription)
	if err != nil {
		return err
	}

	storeclient.NewRetrieveRequest(sourceImage, uploader)
	// TODO check for result
	return nil
}
