package address

// Iterator provides sequential access to address segments with type awareness.
type Iterator struct {
	Address *Address
	index   int
	onKey   bool
}

// NewIterator creates a new iterator for the given address.
func NewIterator(address *Address) *Iterator {
	return &Iterator{
		Address: address,
	}
}

// SegmentType indicates the type of a segment in an address.
type SegmentType uint8

const (
	// SegmentTypeNone indicates no segment (iterator exhausted).
	SegmentTypeNone SegmentType = iota
	// SegmentTypeString indicates a string segment.
	SegmentTypeString
	// SegmentTypeInt indicates an integer index segment.
	SegmentTypeInt
	// SegmentTypeWildcard indicates a wildcard segment (* or #).
	SegmentTypeWildcard
)

// Remainder returns the remaining unvisited portion of the address.
func (i *Iterator) Remainder() *Address {
	return i.Address.From(i.index)
}

// PeekType returns the type of the current segment without advancing.
func (i *Iterator) PeekType() SegmentType {
	if i.index >= len(i.Address.proto.Segments) {
		return SegmentTypeNone
	}
	segment := i.Address.proto.Segments[i.index]
	if i.onKey && segment.IndexInt != nil {
		return SegmentTypeInt
	}
	if segment.Value == "" && segment.IndexInt != nil {
		return SegmentTypeInt
	}
	if segment.Value == "#" || segment.Value == "*" {
		return SegmentTypeWildcard
	}
	return SegmentTypeString
}

// IsDone returns true if the iterator has visited all segments.
func (i *Iterator) IsDone() bool {
	return i.Address == nil || i.Address.proto == nil || i.index >= len(i.Address.proto.Segments)
}

// PeekKey returns the current key without advancing. Returns string key and/or int index.
func (i *Iterator) PeekKey() (*string, *int64) {
	if i.index >= len(i.Address.proto.Segments) {
		return nil, nil
	}
	// use a pointer to avoid large allocations here - v. important as called a LOT
	segment := i.Address.proto.Segments[i.index]
	if i.onKey {
		if k := segment.IndexInt; k != nil {
			return nil, k
		}
		if s := segment.IndexString; s != nil {
			return s, nil
		}
		// should not happen...
		return nil, nil
	}

	if segment.IndexInt != nil || segment.IndexString != nil {
		if segment.Value == "" {
			return segment.IndexString, segment.IndexInt
		}
		return &segment.Value, nil
	}

	return &segment.Value, nil
}

// PopKey returns the current key and advances the iterator. Returns string key and/or int index.
func (i *Iterator) PopKey() (*string, *int64) {
	if i.index >= len(i.Address.proto.Segments) {
		return nil, nil
	}
	// use a pointer to avoid large allocations here - v. important as called a LOT
	segment := i.Address.proto.Segments[i.index]
	if i.onKey {
		if k := segment.IndexInt; k != nil {
			i.index++
			i.onKey = false
			return nil, k
		}
		if s := segment.IndexString; s != nil {
			i.index++
			i.onKey = false
			return s, nil
		}
		// should not happen...
		return nil, nil
	}

	if segment.IndexInt != nil || segment.IndexString != nil {
		if segment.Value == "" {
			i.index++
			return segment.IndexString, segment.IndexInt
		}
		i.onKey = true
		return &segment.Value, nil
	}

	i.index++
	return &segment.Value, nil
}
