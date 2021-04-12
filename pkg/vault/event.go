package vault

import "time"

type Event interface {
	Version() int
	ID() string
	At() time.Time
	Mutate(Entry) (Entry, error)
}

type EventModel struct {
	version int
	id      string
	at      time.Time
}

func (m EventModel) Version() int {
	return m.version
}

func (m EventModel) ID() string {
	return m.id
}

func (m EventModel) At() time.Time {
	return m.at
}

func (m EventModel) Equal(m1 EventModel) bool {
	return m.version == m1.version && m.id == m1.id && m.at.Equal(m1.at)
}

func NewEventModel(version int, id string, at time.Time) EventModel {
	return EventModel{version, id, at}
}
