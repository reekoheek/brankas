package vault

type Event interface {
	Version() int
	ID() string
	Mutate(Entry) (Entry, error)
}

type EventModel struct {
	version int
	id      string
}

func (m EventModel) Version() int {
	return m.version
}

func (m EventModel) ID() string {
	return m.id
}

func NewEventModel(version int, id string) EventModel {
	return EventModel{version, id}
}
