package appstream

type AppStream struct {
	Fleets        []Fleet        `tree:"fleets"`
	ImageBuilders []ImageBuilder `tree:"image_builders"`
}
