package keyvault

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Vault struct {
	resource.Resource `tree:"-"`
	SKUName           value.Value[VaultSKU] `tree:"sku_name"`
}
