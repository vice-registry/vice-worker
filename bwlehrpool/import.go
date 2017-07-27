package bwlehrpool

import (
	"log"
	"github.com/vice-registry/vice-api/models"
	"github.com/vice-registry/vice-worker/storage"
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
	image.Specifics = bwlpImageDetails

	// trigger download of image
	bwlpImageReader, err := handler.GetImageData(bwlpImageID)
	if err != nil {
		log.Printf("Failed to initialise download: %s\n", err)
		return err
	}

	// provide the reader to the storage layer
	err = storage.StoreImage(image, bwlpImageReader)
	if err != nil {
		log.Printf("Failed to provide image reader to storage layer: ", err)
		return err
	}
	log.Printf("Finished import of imageID %s from bwLehrpool", image.ID)
	return nil
}
