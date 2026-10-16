package appflow

type AppFlow struct {
	Flows []Flow `tree:"flows"`
}
