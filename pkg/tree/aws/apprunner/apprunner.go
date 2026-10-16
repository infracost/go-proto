package apprunner

type AppRunner struct {
	Services []Service `tree:"services"`
}
