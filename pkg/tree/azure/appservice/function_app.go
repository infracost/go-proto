package appservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FunctionApp struct {
	resource.Resource  `tree:"-"`
	SKU                value.String                `tree:"sku"`
	Tier               value.Value[AppServiceTier] `tree:"tier"`
	OSType             value.String                `tree:"os_type"`
	MinTLSVersion      value.String                `tree:"min_tls_version"`
	HTTPSOnly          value.Bool                  `tree:"https_only"`
	AppServicePlanID   value.String                `tree:"app_service_plan_id"`
	StorageAccountName value.String                `tree:"storage_account_name"`
	// RuntimeName is the Functions language worker, spelled the way
	// FUNCTIONS_WORKER_RUNTIME spells it: "dotnet-isolated", "dotnet" (the
	// .NET in-process model), "java", "node", "python", "powershell" or
	// "custom". Empty when the IaC does not say.
	RuntimeName value.String `tree:"runtime_name"`
	// RuntimeVersion is the language version for RuntimeName: the major
	// version for .NET, Java and Node.js ("8", "17", "22"), major.minor for
	// Python and PowerShell ("3.11", "7.4").
	RuntimeVersion value.String `tree:"runtime_version"`
}
