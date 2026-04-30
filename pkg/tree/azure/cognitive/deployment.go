package cognitive

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Deployment struct {
	resource.Resource  `tree:"-"`
	Format             value.String `tree:"format"`
	SKU                value.String `tree:"sku"`
	Model              value.String `tree:"model"`
	Version            value.String `tree:"version"`
	Capacity           value.Int    `tree:"capacity"`
	CognitiveAccountID value.String `tree:"cognitive_account_id"`

	Relationships DeploymentRelationships `tree:"-"`
}

type DeploymentRelationships struct {
	Account *Account
}
