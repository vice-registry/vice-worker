package bwlehrpool

import (
	"log"
	"fmt"
	"encoding/json"

	"github.com/OpenSLX/bwlp-go-client/client"
	"github.com/vice-registry/vice-util/models"
	"github.com/vice-registry/vice-util/storeclient"
)

func handleImport(image *models.Image) error {
	log.Printf("Going to import imageID %s from bwLehrpool", image.ID)

	bwlpImageID := image.EnvironmentReference
	log.Printf("bwLehrpool ID: %s\n", bwlpImageID)

	handler, err := Login(image.OriginEnvironment)
	if err != nil {
		log.Printf("Failed to log in to bwLehrpool: ", err)
		return err
	}
	// get image metadata from bwlp-go-client
	bwlpImageDetails, err := handler.GetImageDetails(bwlpImageID)
	if err != nil {
		log.Printf("Failed to retrieve image details: %s\n", err)
		return err
	}

	// store bwlp image and version details as image specifics
	specifics := &bwlpSpecifics{
			Details: bwlpImageDetails,
			Version: nil,
			MachineDescription: nil,
	}
	// get latest version from the list of versions
	version := client.GetLatestVersionDetails(bwlpImageDetails)
	if version == nil {
		return fmt.Errorf("Failed to determine latest version from: %s\n", bwlpImageDetails)
	}
	specifics.Version = version

	// start download process by getting an image downloader (providing a reader)
	bwlpDownloader, err := handler.GetImageData(bwlpImageID)
	if err != nil {
		log.Printf("Failed to initialise download: %s\n", err)
		return err
	}
	specifics.MachineDescription = bwlpDownloader.Ti.MachineDescription

	// encode retrieved informations as json to save as image specifics
	specificsJson, err := json.Marshal(specifics)
	if err != nil {
		log.Printf("Failed to encode bwlp image specifics as json: %s\n", err)
		return err
	}
	image.Specifics = string(specificsJson)

	// provide the reader to the storage layer
	err = storeclient.NewStoreRequest(image, bwlpDownloader)
	if err != nil {
		log.Printf("Failed to request storage of image: ", err)
		return err
	}
	log.Printf("Finished import of imageID %s from bwLehrpool", image.ID)
	return nil
}
