package openstack

import (
	"bufio"
	"log"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/imageservice/v2/images"
	"github.com/vice-registry/vice-util/models"
	"github.com/vice-registry/vice-util/persistence"
	"github.com/vice-registry/vice-util/storeclient"
)

func handleExport(deployment *models.Deployment) error {
	// start export
	log.Printf("Going to export deploymentID %s to OpenStack", deployment.ID)

	// get target environment
	targetEnvironment, err := persistence.GetEnvironment(deployment.EnvironmentID)
	if err != nil {
		return err
	}

	// get image
	image, err := persistence.GetImage(deployment.Imageid)
	if err != nil {
		return err
	}

	// login to openstack
	session, err := login(targetEnvironment)
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

	// upload image to glance
	// TODO read opts from image metadata!
	createOpts := &images.CreateOpts{
		Name:             image.Title,
		ContainerFormat:  "bare",
		DiskFormat:       "qcow2",
		MinDiskGigabytes: 1,
		MinRAMMegabytes:  512,
		Protected:        false,
		//Visibility:       images.ImageVisibilityPrivate,
	}
	newOSImage, err := images.Create(osImageService, createOpts).Extract()

	// create temp file for image caching
	filepath := "/tmp/imagefile"
	file, err := os.Create(filepath)
	if err != nil {
		log.Printf("failed to create imagefile for caching: %s", err)
		return err
	}
	writer := bufio.NewWriter(file)

	// retrieve image from vice-store
	storeclient.NewRetrieveRequest(image, writer)

	// flush to cache file
	if err = writer.Flush(); err != nil {
		log.Printf("failed to flush imagefile for caching: %s", err)
		return err
	}
	file.Close()

	// send data to openstack image service
	file, err = os.Open(filepath)
	result := images.Upload(osImageService, newOSImage.ID, file)
	log.Printf("upload result: %+v, %+v", result.Body, result.Err)

	return nil
}
