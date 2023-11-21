package status

type (
	TopicStatus string
	TableStatus string
)

type TopicInfo struct {
	Errors    map[string]error
	StatusMap map[string]TopicStatus
}

type TableState struct {
	SchemaState string
	TopicState  string
}

const (
	TopicStatusUnknown        TopicStatus = ""
	TopicStatusCreationFailed TopicStatus = "creation failed"
	TopicStatusExists         TopicStatus = "exists"
	TopicStatusCreated        TopicStatus = "created"
	TopicStatusDeleted        TopicStatus = "deleted"

	TableStatusUnkown  TableStatus = ""
	TableStatusUpdated TableStatus = "updated"
	TableStatusCreated TableStatus = "created"
	TableStatusDeleted TableStatus = "deleted"
)

func (tr *TopicInfo) Status(topic string) TopicStatus {
	if tr == nil || tr.StatusMap == nil {
		return TopicStatusUnknown
	}
	if v, ok := tr.StatusMap[topic]; ok {
		return v
	}
	return TopicStatusUnknown
}

func (tr *TopicInfo) Error(topic string) error {
	if v, ok := tr.Errors[topic]; ok {
		return v
	}
	return nil
}
