package openstack

import (
	"io"
	"log"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/imageservice/v2/images"

	"github.com/vice-registry/vice-api/models"
	"github.com/vice-registry/vice-worker/storage"
)

func handleImport(image *models.Image) error {
	// start import
	log.Printf("Going to import imageID %s from OpenStack", image.ID)
	osImageID := image.EnvironmentReference

	// login to openstack
	session, err := login(image.OriginEnvironment)
	if err != nil {
		log.Printf("Unable to log into openstack: %s", err)
		return err
	}

	// get glance image service
	osImageService, err := openstack.NewImageServiceV2(session.Client, gophercloud.EndpointOpts{
		Region: session.Region,
	})
	if err != nil {
		log.Printf("Unable to get image service: %s", err)
		return err
	}

	// read metadata from openstack
	getImageMetadata(osImageService, osImageID)

	// get image reader and provide to storage layer
	reader, err := getImageReader(osImageService, osImageID)
	err = storage.StoreImage(image, reader)

	// close import
	log.Printf("Finished to import imageID %s from OpenStack", image.ID)
	return nil
}

func getImageMetadata(osImageService *gophercloud.ServiceClient, osImageID string) (*images.Image, error) {
	imageGetResult := images.Get(osImageService, osImageID)
	osImage, err := imageGetResult.Extract()
	if err != nil {
		log.Printf("Cannot extract image from imageResult: %s", err)
		return nil, err
	}
	return osImage, nil
}

func getImageReader(osImageService *gophercloud.ServiceClient, osImageID string) (io.Reader, error) {
	getImageDataResult := images.Download(osImageService, osImageID)
	ioReader, err := getImageDataResult.Extract()
	return ioReader, err
}
