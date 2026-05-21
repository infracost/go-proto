package storage

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Queue struct {
	resource.Resource      `tree:"-"`
	StorageAccountName     value.String                        `tree:"storage_account_name"`
	AccountReplicationType value.Value[AccountReplicationType] `tree:"account_replication_type"`
	AccountKind            value.Value[AccountKind]            `tree:"account_kind"`
	AccountTier            value.Value[AccountTier]            `tree:"account_tier"`
}
