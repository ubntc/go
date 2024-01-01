package pebble

type OffsetStatus int

const (
	OffsetStatusOlder OffsetStatus = iota
	OffsetStatusCurrent
	OffsetStatusNewer
)

var OffsetStatuses = []OffsetStatus{OffsetStatusOlder, OffsetStatusCurrent, OffsetStatusNewer}

func (s OffsetStatus) String() string {
	switch s {
	case OffsetStatusOlder:
		return "older"
	case OffsetStatusCurrent:
		return "current"
	case OffsetStatusNewer:
		return "newer"
	default:
		return "unknown"
	}
}

// compareOffset compares the offset of the message with the offset of the
// last message seen. NOTE: Must be protected by r.mu!
func CompareOffsetByKey(currentKey, otherKey []byte) OffsetStatus {
	if currentKey == nil {
		// no message seen so far
		return OffsetStatusNewer
	}
	if otherKey == nil {
		// no new message
		return OffsetStatusOlder
	}
	currentOffset := Offset(currentKey)
	offset := Offset(otherKey)
	switch {
	case offset < currentOffset:
		return OffsetStatusOlder
	case offset == currentOffset:
		return OffsetStatusCurrent
	default:
		return OffsetStatusNewer
	}
}
