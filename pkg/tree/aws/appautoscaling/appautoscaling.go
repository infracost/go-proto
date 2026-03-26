package appautoscaling

type AppAutoScaling struct {
	Targets []Target `tree:"targets"`
}
