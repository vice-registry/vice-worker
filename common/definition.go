package adaptors

import "omi-gitlab.e-technik.uni-ulm.de/vice/vice-api/models"

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
