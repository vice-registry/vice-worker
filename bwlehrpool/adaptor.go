package bwlehrpool

import (
	"github.com/vice-registry/vice-util/models"
)

// BwLehrpool defines the adaptor for bwLehrpool
type BwLehrpool struct {
}

// Import functionality of the bwLehrpool adaptor
func (adaptor BwLehrpool) Import(image *models.Image) error {
	return handleImport(image)
}

// Deploy functionality of the bwLehrpool adaptor
func (adaptor BwLehrpool) Deploy(deployment *models.Deployment) error {
	return handleExport(deployment)
}
