package cloudwatch

type CloudWatch struct {
	Dashboards   []Dashboard   `tree:"dashboards"`
	EventBuses   []EventBus    `tree:"event_buses"`
	LogGroups    []LogGroup    `tree:"log_groups"`
	MetricAlarms []MetricAlarm `tree:"metric_alarms"`
	EventTargets []EventTarget `tree:"event_targets"`
}
