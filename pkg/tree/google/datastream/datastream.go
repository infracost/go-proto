package datastream

type Datastream struct {
	Streams []Stream `tree:"streams"`
}
