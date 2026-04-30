package storage

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ManagementPolicy struct {
	resource.Resource `tree:"-"`
	StorageAccountID  value.String `tree:"storage_account_id"`
}
