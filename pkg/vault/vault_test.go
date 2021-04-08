package vault

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	table := []struct {
		name    string
		aEvents []Event
	}{
		{
			"positive case #1",
			nil,
		},
		{
			"positive case #2",
			[]Event{
				tNewEvent(1, "", nil),
				tNewEvent(2, "", nil),
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.aEvents)

			assert.Equal(t, 0, len(v.uncommitedEvents))
			assert.Equal(t, len(tt.aEvents), len(v.commitedEvents))
			assert.Equal(t, tt.aEvents, v.commitedEvents)
		})
	}
}

func TestVault_Apply(t *testing.T) {
	table := []struct {
		name   string
		vault  Vault
		aEvent tEvent
		rErr   string
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
				assert.Equal(t, tt.rErr, err.Error())
				return
			}
			assert.Equal(t, "", tt.rErr, "Expecting error: %s", tt.rErr)
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
