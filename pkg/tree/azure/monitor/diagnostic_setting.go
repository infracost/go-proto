package monitor

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DiagnosticSetting struct {
	resource.Resource      `tree:"-"`
	EventHubTarget         value.Bool   `tree:"event_hub_target"`
	PartnerSolutionTarget  value.Bool   `tree:"partner_solution_target"`
	StorageAccountTarget   value.Bool   `tree:"storage_account_target"`
	TargetResourceID       value.String `tree:"target_resource_id"`
}
