package automation

type Automation struct {
	Accounts              []Account              `tree:"accounts"`
	JobSchedules          []JobSchedule          `tree:"job_schedules"`
	DSCConfigurations     []DSCConfiguration     `tree:"dsc_configurations"`
	DSCNodeConfigurations []DSCNodeConfiguration `tree:"dsc_node_configurations"`
}
