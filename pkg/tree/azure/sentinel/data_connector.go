package sentinel

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DataConnector struct {
	resource.Resource `tree:"-"`
	WorkspaceID       value.String `tree:"workspace_id"`
}
