package cloudtrail

type CloudTrail struct {
	Trails []Trail `tree:"trails"`
}
