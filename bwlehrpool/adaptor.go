package bwlehrpool

import (
	"fmt"

	"github.com/vice-registry/vice-api/models"
)

// BwLehrpool defines the adaptor for OpenStack
type BwLehrpool struct {
}

// Import functionality of the bwLehrpool adaptor
func (adaptor BwLehrpool) Import(image *models.Image) error {
	err := fmt.Errorf("no implementation for adaptor bwLehrpool")
	return err
}

// Deploy functionality of the bwLehrpool adaptor
func (adaptor BwLehrpool) Deploy(deployment *models.Deployment) error {
	err := fmt.Errorf("no implementation for adaptor bwLehrpool")
	return err
}
