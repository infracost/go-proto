package datafactory

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type IntegrationRuntimeAzureSSIS struct {
	resource.Resource `tree:"-"`
	NumberOfNodes     value.Int                    `tree:"number_of_nodes"`
	NodeSize          value.String                 `tree:"node_size"`
	Edition           value.Value[SSISEdition]     `tree:"edition"`
	LicenseType       value.Value[SSISLicenseType] `tree:"license_type"`
}
