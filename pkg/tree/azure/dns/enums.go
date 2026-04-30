package dns

type ZoneType uint32

const (
	ZoneTypeUnknown ZoneType = iota
	ZoneTypePublic
	ZoneTypePrivate
)
