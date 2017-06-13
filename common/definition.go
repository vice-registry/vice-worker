package adaptors

import "github.com/vice-registry/vice-api/models"

const (
	// WorkerTypeImport captures enum value "import"
	WorkerTypeImport string = "import"
	// WorkerTypeExport captures enum value "export"
	WorkerTypeExport string = "export"
)

// WorkerType defines type of worker instance (WorkerTypeImport or WorkerTypeExport)
var WorkerType string

// Adaptor interface for environment adaptors
type Adaptor interface {
	Import(image *models.Image) error
	Deploy(deployment *models.Deployment) error
}
