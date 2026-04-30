package synapse

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type SQLPool struct {
	resource.Resource      `tree:"-"`
	SKUName                value.String `tree:"sku_name"`
	WorkspaceID            value.String `tree:"workspace_id"`
	GeoBackupPolicyEnabled value.Bool   `tree:"geo_backup_policy_enabled"`

	Relationships SQLPoolRelationships `tree:"-"`
}

type SQLPoolRelationships struct {
	Workspace *Workspace
}
