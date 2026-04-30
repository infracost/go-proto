package scheduler

type Scheduler struct {
	Schedules []Schedule `tree:"schedules"`
}
