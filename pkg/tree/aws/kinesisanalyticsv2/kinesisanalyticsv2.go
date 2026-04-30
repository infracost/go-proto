package kinesisanalyticsv2

type KinesisAnalyticsV2 struct {
	Applications         []Application         `tree:"applications"`
	ApplicationSnapshots []ApplicationSnapshot `tree:"application_snapshots"`
}
