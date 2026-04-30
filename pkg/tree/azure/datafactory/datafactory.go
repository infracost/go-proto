package datafactory

type DataFactory struct {
	DataFactories                []Factory                    `tree:"data_factories"`
	IntegrationRuntimeAzures     []IntegrationRuntimeAzure     `tree:"integration_runtime_azures"`
	IntegrationRuntimeSelfHosteds []IntegrationRuntimeSelfHosted `tree:"integration_runtime_self_hosteds"`
	IntegrationRuntimeManageds   []IntegrationRuntimeManaged   `tree:"integration_runtime_manageds"`
	IntegrationRuntimeAzureSSISs []IntegrationRuntimeAzureSSIS `tree:"integration_runtime_azure_ssiss"`
}
