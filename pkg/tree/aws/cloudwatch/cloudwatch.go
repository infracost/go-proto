package cloudwatch

type CloudWatch struct {
	Dashboards              []Dashboard              `tree:"dashboards"`
	EventBuses              []EventBus               `tree:"event_buses"`
	LogGroups               []LogGroup               `tree:"log_groups"`
	MetricAlarms            []MetricAlarm            `tree:"metric_alarms"`
	EventTargets            []EventTarget            `tree:"event_targets"`
	CompositeAlarms         []CompositeAlarm         `tree:"composite_alarms"`
	ContributorInsightRules []ContributorInsightRule `tree:"contributor_insight_rules"`
	MetricStreams           []MetricStream           `tree:"metric_streams"`
	SyntheticsCanaries      []SyntheticsCanary       `tree:"synthetics_canaries"`
	EventArchives           []EventArchive           `tree:"event_archives"`
	EventAPIDestinations    []EventAPIDestination    `tree:"event_api_destinations"`
	InternetMonitors        []InternetMonitor        `tree:"internet_monitors"`
	RUMAppMonitors          []RUMAppMonitor          `tree:"rum_app_monitors"`
}
