package azure

import (
	"github.com/infracost/go-proto/pkg/tree/azure/activedirectory"
	"github.com/infracost/go-proto/pkg/tree/azure/apimanagement"
	"github.com/infracost/go-proto/pkg/tree/azure/appconfiguration"
	"github.com/infracost/go-proto/pkg/tree/azure/applicationinsights"
	"github.com/infracost/go-proto/pkg/tree/azure/appservice"
	"github.com/infracost/go-proto/pkg/tree/azure/automation"
	"github.com/infracost/go-proto/pkg/tree/azure/cognitive"
	"github.com/infracost/go-proto/pkg/tree/azure/compute"
	"github.com/infracost/go-proto/pkg/tree/azure/container"
	"github.com/infracost/go-proto/pkg/tree/azure/cosmosdb"
	"github.com/infracost/go-proto/pkg/tree/azure/database"
	"github.com/infracost/go-proto/pkg/tree/azure/databricks"
	"github.com/infracost/go-proto/pkg/tree/azure/datafactory"
	"github.com/infracost/go-proto/pkg/tree/azure/dns"
	"github.com/infracost/go-proto/pkg/tree/azure/eventgrid"
	"github.com/infracost/go-proto/pkg/tree/azure/eventhub"
	"github.com/infracost/go-proto/pkg/tree/azure/expressroute"
	"github.com/infracost/go-proto/pkg/tree/azure/fabric"
	"github.com/infracost/go-proto/pkg/tree/azure/firewall"
	"github.com/infracost/go-proto/pkg/tree/azure/frontdoor"
	"github.com/infracost/go-proto/pkg/tree/azure/hdinsight"
	"github.com/infracost/go-proto/pkg/tree/azure/identityplatform"
	"github.com/infracost/go-proto/pkg/tree/azure/integration"
	"github.com/infracost/go-proto/pkg/tree/azure/iothub"
	"github.com/infracost/go-proto/pkg/tree/azure/keyvault"
	"github.com/infracost/go-proto/pkg/tree/azure/loganalytics"
	"github.com/infracost/go-proto/pkg/tree/azure/logicapps"
	"github.com/infracost/go-proto/pkg/tree/azure/machinelearning"
	"github.com/infracost/go-proto/pkg/tree/azure/monitor"
	"github.com/infracost/go-proto/pkg/tree/azure/network"
	"github.com/infracost/go-proto/pkg/tree/azure/notificationhub"
	"github.com/infracost/go-proto/pkg/tree/azure/powerbi"
	"github.com/infracost/go-proto/pkg/tree/azure/recoveryservices"
	"github.com/infracost/go-proto/pkg/tree/azure/redis"
	"github.com/infracost/go-proto/pkg/tree/azure/search"
	"github.com/infracost/go-proto/pkg/tree/azure/securitycenter"
	"github.com/infracost/go-proto/pkg/tree/azure/sentinel"
	"github.com/infracost/go-proto/pkg/tree/azure/servicebus"
	"github.com/infracost/go-proto/pkg/tree/azure/signalr"
	"github.com/infracost/go-proto/pkg/tree/azure/storage"
	"github.com/infracost/go-proto/pkg/tree/azure/synapse"
)

type Azure struct {
	ActiveDirectory     activedirectory.ActiveDirectory       `tree:"activedirectory"`
	APIManagement       apimanagement.APIManagement           `tree:"apimanagement"`
	AppConfiguration    appconfiguration.AppConfiguration     `tree:"appconfiguration"`
	ApplicationInsights applicationinsights.ApplicationInsights `tree:"applicationinsights"`
	AppService          appservice.AppService                 `tree:"appservice"`
	Automation          automation.Automation                 `tree:"automation"`
	Cognitive           cognitive.Cognitive                   `tree:"cognitive"`
	Compute             compute.Compute                       `tree:"compute"`
	Container           container.Container                   `tree:"container"`
	CosmosDB            cosmosdb.CosmosDB                     `tree:"cosmosdb"`
	Database            database.Database                     `tree:"database"`
	Databricks          databricks.Databricks                 `tree:"databricks"`
	DataFactory         datafactory.DataFactory               `tree:"datafactory"`
	DNS                 dns.DNS                               `tree:"dns"`
	EventGrid           eventgrid.EventGrid                   `tree:"eventgrid"`
	EventHub            eventhub.EventHub                     `tree:"eventhub"`
	ExpressRoute        expressroute.ExpressRoute             `tree:"expressroute"`
	Fabric              fabric.Fabric                         `tree:"fabric"`
	Firewall            firewall.Firewall                     `tree:"firewall"`
	FrontDoor           frontdoor.FrontDoor                   `tree:"frontdoor"`
	HDInsight           hdinsight.HDInsight                   `tree:"hdinsight"`
	IdentityPlatform    identityplatform.IdentityPlatform     `tree:"identityplatform"`
	Integration         integration.Integration               `tree:"integration"`
	IoTHub              iothub.IoTHub                         `tree:"iothub"`
	KeyVault            keyvault.KeyVault                     `tree:"keyvault"`
	LogAnalytics        loganalytics.LogAnalytics             `tree:"loganalytics"`
	LogicApps           logicapps.LogicApps                   `tree:"logicapps"`
	MachineLearning     machinelearning.MachineLearning       `tree:"machinelearning"`
	Monitor             monitor.Monitor                       `tree:"monitor"`
	Network             network.Network                       `tree:"network"`
	NotificationHub     notificationhub.NotificationHub       `tree:"notificationhub"`
	PowerBI             powerbi.PowerBI                       `tree:"powerbi"`
	RecoveryServices    recoveryservices.RecoveryServices     `tree:"recoveryservices"`
	Redis               redis.Redis                           `tree:"redis"`
	Search              search.Search                         `tree:"search"`
	SecurityCenter      securitycenter.SecurityCenter         `tree:"securitycenter"`
	Sentinel            sentinel.Sentinel                     `tree:"sentinel"`
	ServiceBus          servicebus.ServiceBus                 `tree:"servicebus"`
	SignalR             signalr.SignalR                       `tree:"signalr"`
	Storage             storage.Storage                       `tree:"storage"`
	Synapse             synapse.Synapse                       `tree:"synapse"`
}

func (az *Azure) PostProcess() {
	// NOTE: Service-level PostProcess() methods are invoked automatically by the
	// reflective tree walker in tree.go. Only cross-service wiring belongs here.

	// cross-service: link backup protected VMs to compute virtual machines
	for i, vm := range az.RecoveryServices.BackupProtectedVMs {
		for j := range az.Compute.VirtualMachines {
			if vm.SourceVMID.Value() == az.Compute.VirtualMachines[j].ID {
				az.RecoveryServices.BackupProtectedVMs[i].Relationships.SourceVM = &az.Compute.VirtualMachines[j]
				break
			}
		}
	}

	// cross-service: link logic app standards to app service plans
	for i, standard := range az.LogicApps.Standards {
		if standard.AppServicePlanID.IsEmpty() {
			continue
		}
		for _, plan := range az.AppService.AppServicePlans {
			if plan.ID == standard.AppServicePlanID.Value() {
				az.LogicApps.Standards[i].SKU = plan.SKUSize
				break
			}
		}
		for _, plan := range az.AppService.ServicePlans {
			if plan.ID == standard.AppServicePlanID.Value() {
				az.LogicApps.Standards[i].SKU = plan.SKUName
				break
			}
		}
	}

	// cross-service: link function apps to storage accounts
	for _, fa := range az.AppService.FunctionApps {
		for j := range az.Storage.Accounts {
			if fa.StorageAccountName.Value() == az.Storage.Accounts[j].Name.Value() {
				az.Storage.Accounts[j].UsedByFunctionApps = true
				break
			}
		}
	}
}
