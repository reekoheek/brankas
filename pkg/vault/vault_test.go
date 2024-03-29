package vault

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
)

func TestNew(t *testing.T) {
	table := []struct {
		name    string
		aID     string
		aEvents []Event
	}{
		{
			"positive case #1",
			"foo",
			nil,
		},
		{
			"positive case #2",
			"bar",
			[]Event{
				tNewEvent(1, "", nil),
				tNewEvent(2, "", nil),
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.aID, tt.aEvents)

			assert.Equal(t, tt.aID, v.id)
			assert.Equal(t, 0, len(v.uncommitedEvents))
			assert.Equal(t, len(tt.aEvents), len(v.commitedEvents))
			assert.DeepEqual(t, tt.aEvents, v.commitedEvents, cmp.Comparer(func(ev1 Event, ev2 Event) bool {
				return ev1.ID() == ev2.ID() && ev1.Version() == ev2.Version()
			}))
		})
	}
}

func TestVault_Apply(t *testing.T) {
	table := []struct {
		name   string
		vault  Vault
		aEvent tEvent
		xErr   string
	}{
		{
			"positive case",
			Vault{entries: map[string]Entry{}},
			tNewEvent(1, "", nil),
			"",
		},
		{
			"invalid id",
			Vault{
				entries:        map[string]Entry{},
				commitedEvents: []Event{tNewEvent(10, "", nil)},
			},
			tNewEvent(1, "", nil),
			"invalid event version",
		},
		{
			"err on event mutate",
			Vault{
				entries:        map[string]Entry{},
				commitedEvents: []Event{tNewEvent(10, "", nil)},
			},
			tNewEvent(11, "", func(e Entry) (Entry, error) { return e, fmt.Errorf("mutate err") }),
			"mutate err",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vault.Apply(tt.aEvent)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}
			assert.Equal(t, "", tt.xErr, "Expecting error: %s", tt.xErr)
			assert.Equal(t, 1, len(tt.vault.uncommitedEvents))
		})
	}
}

func TestVault_Version(t *testing.T) {
	table := []struct {
		name     string
		vault    Vault
		xVersion int
	}{
		{
			"empty",
			Vault{},
			0,
		},
		{
			"commited only",
			Vault{commitedEvents: []Event{tNewEvent(3, "", nil), tNewEvent(4, "", nil)}},
			4,
		},
		{
			"dirty",
			Vault{uncommitedEvents: []Event{tNewEvent(3, "", nil), tNewEvent(4, "", nil)}},
			4,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.xVersion, tt.vault.Version())
		})
	}
}

func TestVault_Events(t *testing.T) {
	table := []struct {
		name     string
		v        Vault
		aVersion int
		xLen     int
	}{
		{
			"1",
			Vault{
				commitedEvents: []Event{tNewEvent(1, "", nil), tNewEvent(2, "", nil)},
			},
			1,
			2,
		},
		{
			"2",
			Vault{
				commitedEvents: []Event{tNewEvent(1, "", nil), tNewEvent(2, "", nil)},
			},
			2,
			1,
		},
		{
			"2",
			Vault{},
			1,
			0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			evs := tt.v.Events(tt.aVersion)
			assert.Equal(t, tt.xLen, len(evs))
		})
	}
}

func TestVault_UncommitedEvents(t *testing.T) {
	table := []struct {
		name string
		v    Vault
		xLen int
	}{
		{
			"1",
			Vault{
				uncommitedEvents: []Event{tNewEvent(1, "", nil), tNewEvent(2, "", nil)},
			},
			2,
		},
		{
			"2",
			Vault{},
			0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			evs := tt.v.UncommitedEvents()
			assert.Equal(t, tt.xLen, len(evs))
		})
	}
}

type tEvent struct {
	EventModel
	fn func(Entry) (Entry, error)
}

func (t tEvent) Mutate(e Entry) (Entry, error) {
	if t.fn == nil {
		return e, nil
	}

	return t.fn(e)
}

func tNewEvent(version int, id string, fn func(Entry) (Entry, error)) tEvent {
	return tEvent{
		EventModel: EventModel{
			version: version,
			id:      id,
		},
		fn: fn,
	}
}
