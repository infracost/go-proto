package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infracost/go-proto/pkg/tree/azure/appservice"
	"github.com/infracost/go-proto/pkg/tree/azure/compute"
	"github.com/infracost/go-proto/pkg/tree/azure/loganalytics"
	"github.com/infracost/go-proto/pkg/tree/azure/logicapps"
	"github.com/infracost/go-proto/pkg/tree/azure/recoveryservices"
	"github.com/infracost/go-proto/pkg/tree/azure/sentinel"
	"github.com/infracost/go-proto/pkg/tree/azure/storage"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

func TestAzurePostProcess_IsIdempotent(t *testing.T) {
	az := &Azure{
		LogAnalytics: loganalytics.LogAnalytics{
			Workspaces: []loganalytics.Workspace{
				{
					Resource:        resource.Resource{ID: "ws-1"},
					SentinelEnabled: value.New(false, 0, "", nil),
				},
			},
		},
		Sentinel: sentinel.Sentinel{
			DataConnectors: []sentinel.DataConnector{
				{
					Resource:    resource.Resource{ID: "dc-1"},
					WorkspaceID: value.New("ws-1", 0, "", nil),
				},
			},
		},
		RecoveryServices: recoveryservices.RecoveryServices{
			BackupProtectedVMs: []recoveryservices.BackupProtectedVM{
				{
					Resource:   resource.Resource{ID: "bpvm-1"},
					SourceVMID: value.New("vm-1", 0, "", nil),
				},
			},
		},
		Compute: compute.Compute{
			VirtualMachines: []compute.VirtualMachine{
				{Resource: resource.Resource{ID: "vm-1"}},
			},
		},
		LogicApps: logicapps.LogicApps{
			Standards: []logicapps.Standard{
				{
					Resource:         resource.Resource{ID: "la-1"},
					AppServicePlanID: value.New("plan-1", 0, "", nil),
				},
			},
		},
		AppService: appservice.AppService{
			AppServicePlans: []appservice.AppServicePlan{
				{Resource: resource.Resource{ID: "plan-1"}, SKUSize: value.New("S1", 0, "", nil)},
			},
			FunctionApps: []appservice.FunctionApp{
				{
					Resource:           resource.Resource{ID: "fa-1"},
					StorageAccountName: value.New("sa-1", 0, "", nil),
				},
			},
		},
		Storage: storage.Storage{
			Accounts: []storage.Account{
				{
					Resource:           resource.Resource{ID: "acc-1"},
					Name:               value.New("sa-1", 0, "", nil),
					UsedByFunctionApps: value.New(false, 0, "", nil),
				},
			},
		},
	}

	az.PostProcess()
	sentinelEnabled := az.LogAnalytics.Workspaces[0].SentinelEnabled.Value()
	bpvmSourceVM := az.RecoveryServices.BackupProtectedVMs[0].Relationships.SourceVM
	standardSKU := az.LogicApps.Standards[0].SKU.Value()
	saUsedByFA := az.Storage.Accounts[0].UsedByFunctionApps.Value()

	az.PostProcess()
	assert.Equal(t, sentinelEnabled, az.LogAnalytics.Workspaces[0].SentinelEnabled.Value())
	assert.Equal(t, bpvmSourceVM, az.RecoveryServices.BackupProtectedVMs[0].Relationships.SourceVM)
	assert.Equal(t, standardSKU, az.LogicApps.Standards[0].SKU.Value())
	assert.Equal(t, saUsedByFA, az.Storage.Accounts[0].UsedByFunctionApps.Value())

	assert.True(t, sentinelEnabled)
	assert.NotNil(t, bpvmSourceVM)
	assert.Equal(t, "S1", standardSKU)
	assert.True(t, saUsedByFA)
}
