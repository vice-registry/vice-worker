package openstack

import (
	"log"

	"omi-gitlab.e-technik.uni-ulm.de/vice/vice-api/persistence"
)

// Test import via OpenStack import adaptor
func Test() {
	image, err := persistence.GetImage("VbhV4v")
	if err != nil {
		log.Fatalf("Cannot get image VbhV4v: %s", err)
	}

	err = Import(image)
	log.Printf("Import: %s", err)
}
