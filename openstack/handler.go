package openstack

import (
	"io"
	"log"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/imageservice/v2/images"

	"fmt"

	"omi-gitlab.e-technik.uni-ulm.de/vice/vice-api/models"
	"omi-gitlab.e-technik.uni-ulm.de/vice/vice-import/storage"
)

// Credentials credential extension for OpenStack
type Credentials struct {
	TenantName string
	Region     string
}

// Import an image from OpenStack
func Import(image *models.Image) error {

	// start import
	log.Printf("Going to import imageID %s from OpenStack", image.ID)

	// extract data from image model
	endpoint := image.OriginEnvironment.Credentials.Location
	username := image.OriginEnvironment.Credentials.Username
	password := image.OriginEnvironment.Credentials.Password
	specifics, ok := image.OriginEnvironment.Credentials.Specifics.(map[string]interface{})
	if !ok {
		err := fmt.Errorf("failed to convert specifics")
		log.Printf("Unable to get openstack specifics: %s", err)
		return err
	}
	rawTenant, ok := specifics["TenantName"]
	if !ok {
		err := fmt.Errorf("missing specific value %s", "TenantName")
		log.Printf("Unable to get openstack specifics: %s", err)
		return err
	}
	tenant := rawTenant.(string)
	rawRegion, ok := specifics["Region"]
	if !ok {
		err := fmt.Errorf("missing specific value %s", "Region")
		log.Printf("Unable to get openstack specifics: %s", err)
		return err
	}
	region := rawRegion.(string)
	osImageID := image.EnvironmentReference

	// login to openstack
	osProvider, err := login(endpoint, username, password, tenant)
	if err != nil {
		log.Printf("Unable to login to openstack: %s", err)
		return err
	}

	// get glance image service
	osImageService, err := openstack.NewImageServiceV2(osProvider, gophercloud.EndpointOpts{
		Region: region,
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

func login(endpoint string, username string, password string, tenant string) (*gophercloud.ProviderClient, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: endpoint,
		Username:         username,
		Password:         password,
		TenantName:       tenant,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Printf("Unable to authenticate at openstack: %s", err)
		return nil, err
	}
	return provider, nil
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
