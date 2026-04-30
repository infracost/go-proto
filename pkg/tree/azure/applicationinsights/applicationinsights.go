package applicationinsights

type ApplicationInsights struct {
	Insights         []Insights         `tree:"insights"`
	StandardWebTests []StandardWebTest `tree:"standard_web_tests"`
	WebTests         []WebTest          `tree:"web_tests"`
}
