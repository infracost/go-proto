package grafana

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Workspace struct {
	resource.Resource `tree:"-"`

	Relationships WorkspaceRelationships `tree:"-"`
}

type WorkspaceRelationships struct {
	LicenseAssociation *LicenseAssociation
}

type LicenseAssociation struct {
	resource.Resource `tree:"-"`
	WorkspaceID       value.String `tree:"workspace_id"`
	Enterprise        value.Bool   `tree:"enterprise"`
}
