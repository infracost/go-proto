package keyvault

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Key struct {
	resource.Resource `tree:"-"`
	KeyType           value.Value[KeyType] `tree:"key_type"`
	KeySize           value.String         `tree:"key_size"`
	SKUName           value.Value[VaultSKU] `tree:"sku_name"`
	KeyVaultID        value.String `tree:"key_vault_id"`
}
