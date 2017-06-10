package actions

import (
	"fmt"
	"log"
	"strings"

	"omi-gitlab.e-technik.uni-ulm.de/vice/vice-api/persistence"
	"omi-gitlab.e-technik.uni-ulm.de/vice/vice-import/bwlehrpool"
	"omi-gitlab.e-technik.uni-ulm.de/vice/vice-import/openstack"
)

func handleAction(imageID string) error {

	image, err := persistence.GetImage(imageID)
	if err != nil {
		log.Printf("Error in import handleAction for imageID %s: %s", imageID, err)
		return err
	}
	log.Printf("handle import for image %+v", image)

	importAdaptor := strings.ToLower(image.OriginEnvironment.Managementlayer.Software)
	switch importAdaptor {
	case "openstack":
		err := openstack.Import(image)
		if err != nil {
			log.Printf("Error in import handleAction: %s", err)
			return err
		}
	case "bwlehrpool":
		err := bwlehrpool.Import(image)
		if err != nil {
			log.Printf("Error in import handleAction: %s", err)
			return err
		}
	default:
		err := fmt.Errorf("no import adaptor found for management layer %s", importAdaptor)
		log.Printf("Error in import handleAction: %s", err)
		return err
	}

	// Update image in Couchbase
	image.Imported = true
	persistence.UpdateImage(image)

	return nil
}
