package openstack

import (
	"fmt"
	"log"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/vice-registry/vice-util/models"
)

// https://github.com/hashicorp/terraform/blob/401c6a95a7d98004b2fa4d0f8d68fc4b95f870e9/builtin/providers/openstack/resource_openstack_images_image_v2.go#L170

// session representaton of user session to openstack
type session struct {
	TenantName string
	Region     string
	Client     *gophercloud.ProviderClient
}

func login(environment *models.Environment) (session, error) {
	// extract data from image model
	endpoint := environment.Credentials.Location
	username := environment.Credentials.Username
	password := environment.Credentials.Password
	specifics, ok := environment.Credentials.Specifics.(map[string]interface{})
	if !ok {
		err := fmt.Errorf("failed to convert specifics")
		log.Printf("Unable to get openstack specifics: %s", err)
		return session{}, err
	}

	// extract tenant
	rawTenant, ok := specifics["TenantName"]
	if !ok {
		err := fmt.Errorf("missing specific value %s", "TenantName")
		log.Printf("Unable to get openstack specifics: %s", err)
		return session{}, err
	}
	tenant := rawTenant.(string)

	// extract region
	rawRegion, ok := specifics["Region"]
	if !ok {
		err := fmt.Errorf("missing specific value %s", "Region")
		log.Printf("Unable to get openstack specifics: %s", err)
		return session{}, err
	}
	region := rawRegion.(string)

	// extract region
	rawDomain, ok := specifics["Domain"]
	if !ok {
		err := fmt.Errorf("missing specific value %s", "Domain")
		log.Printf("Unable to get openstack specifics: %s", err)
		return session{}, err
	}
	domain := rawDomain.(string)

	// login to openstack
	osProvider, err := getProviderClient(endpoint, username, password, tenant, domain)
	if err != nil {
		log.Printf("Unable to login to openstack: %s", err)
		return session{}, err
	}

	// return new session
	session := session{
		Client:     osProvider,
		Region:     region,
		TenantName: tenant,
	}
	return session, nil

}

func getProviderClient(endpoint string, username string, password string, tenant string, domainname string) (*gophercloud.ProviderClient, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: endpoint,
		Username:         username,
		Password:         password,
		TenantName:       tenant,
		DomainName:       domainname,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Printf("Unable to authenticate at openstack: %s", err)
		return nil, err
	}
	return provider, nil
}
