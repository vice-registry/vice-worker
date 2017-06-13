package openstack

import "github.com/vice-registry/vice-api/models"

// OpenStack defines the adaptor for OpenStack
type OpenStack struct {
}

// Import functionality of the openstack adaptor
func (adaptor OpenStack) Import(image *models.Image) error {
	return handleImport(image)
}

// Deploy functionality of the openstack adaptor
func (adaptor OpenStack) Deploy(deployment *models.Deployment) error {
	return handleExport(deployment)
}
