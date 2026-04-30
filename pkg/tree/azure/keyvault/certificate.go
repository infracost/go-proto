package keyvault

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Certificate struct {
	resource.Resource `tree:"-"`
	KeyVaultID        value.String `tree:"key_vault_id"`
}
