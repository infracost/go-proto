package storage

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Account struct {
	resource.Resource            `tree:"-"`
	Name                         value.String `tree:"name"`
	AccessTier                   value.Value[AccessTier]             `tree:"access_tier"`
	AccountKind                  value.Value[AccountKind]            `tree:"account_kind"`
	AccountReplicationType       value.Value[AccountReplicationType] `tree:"account_replication_type"`
	AccountTier                  value.Value[AccountTier]            `tree:"account_tier"`
	NFSv3                        value.Bool   `tree:"nfsv3"`
	LastTimeBlobEnabled          value.Bool   `tree:"last_time_blob_enabled"`
	PublicNetworkAccessEnabled   value.Bool   `tree:"public_network_access_enabled"`
	UsedByFunctionApps           value.Bool   `tree:"used_by_function_apps"`

	Relationships AccountRelationships `tree:"-"`
}

type AccountRelationships struct {
	ManagementPolicy *ManagementPolicy
}
