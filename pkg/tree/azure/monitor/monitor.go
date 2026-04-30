package monitor

type Monitor struct {
	ActionGroups                []ActionGroup                `tree:"action_groups"`
	DiagnosticSettings          []DiagnosticSetting          `tree:"diagnostic_settings"`
	MetricAlerts                []MetricAlert                `tree:"metric_alerts"`
	ScheduledQueryRulesAlerts   []ScheduledQueryRulesAlert   `tree:"scheduled_query_rules_alerts"`
	DataCollectionRules         []DataCollectionRule         `tree:"data_collection_rules"`
}
