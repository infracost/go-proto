package datasync

type DataSync struct {
	Tasks []Task `tree:"tasks"`
}
