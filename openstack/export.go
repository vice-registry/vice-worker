package openstack

import (
	"log"

	"github.com/vice-registry/vice-util/models"
)

func handleExport(deployment *models.Deployment) error {
	// start export
	log.Printf("Going to export deploymentID %s to OpenStack", deployment.ID)

	log.Printf("Temporarily disabled!")

	/*

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
		createOpts := &images.CreateOpts{
			Name:             image.ID, // TODO
			ContainerFormat:  "bare",
			DiskFormat:       "qcow2",
			MinDiskGigabytes: 1,
			MinRAMMegabytes:  512,
			Protected:        false,
			//Visibility:       images.ImageVisibilityPublic,
		}
		newImg, err := images.Create(osImageService, createOpts).Extract()
		imgFile, err := storage.RetrieveImage(image)
		defer imgFile.Close()
		images.Upload(osImageService, newImg.ID, imgFile)

	*/

	return nil
}
