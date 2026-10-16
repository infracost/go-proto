package cloudtrail

type CloudTrail struct {
	Trails          []Trail          `tree:"trails"`
	EventDataStores []EventDataStore `tree:"event_data_stores"`
}
