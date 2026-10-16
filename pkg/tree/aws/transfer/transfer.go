package transfer

type Transfer struct {
	Servers    []Server    `tree:"servers"`
	Connectors []Connector `tree:"connectors"`
	WebApps    []WebApp    `tree:"web_apps"`
}
