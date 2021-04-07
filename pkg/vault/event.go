package vault

type Event interface {
	Version() int
	ID() string
	Apply(Entry) error
}

type EventModel struct {
	version int
	id      string
}

func (m *EventModel) Version() int {
	return m.version
}

func (m *EventModel) ID() string {
	return m.id
}
