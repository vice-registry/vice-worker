package actions

import (
	"log"

	"github.com/vice-registry/vice-api/persistence"
)

func handleImport(action Action) error {
	image := action.image
	log.Printf("handle import for image %+v", image)

	// call adaptor
	err := action.adaptor.Import(image)
	if err != nil {
		log.Printf("Error in handleImport for imageID %s: %s", action.reference, err)
		return err
	}

	// Update image in Couchbase
	image.Imported = true
	_, err = persistence.UpdateImage(image)
	if err != nil {
		log.Printf("Error in handleImport for imageID %s: %s", action.reference, err)
		return err
	}

	return nil
}
