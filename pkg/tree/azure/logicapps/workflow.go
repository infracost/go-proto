package logicapps

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

// Workflow represents an azurerm_logic_app_workflow (Consumption plan) - a
// serverless Logic App with no dedicated compute SKU of its own. Azure bills
// it entirely per Action execution and connector call, so all of its cost
// comes from usage data rather than Terraform config.
type Workflow struct {
	resource.Resource `tree:"-"`
}
