package logicapps

type LogicApps struct {
	Standards           []Standard           `tree:"standards"`
	IntegrationAccounts []IntegrationAccount `tree:"integration_accounts"`
}
