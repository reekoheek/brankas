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
				tEvent{
					EventModel: &EventModel{1, ""},
				},
				tEvent{
					EventModel: &EventModel{2, ""},
				},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.aEvents)

			assert.Equal(t, 0, len(v.uncomittedEvents))
			assert.Equal(t, len(tt.aEvents), len(v.commitedEvents))
			assert.Equal(t, tt.aEvents, v.commitedEvents)
		})
	}
}

func TestVault_Apply(t *testing.T) {
	table := []struct {
		name   string
		Vault  Vault
		aEvent tEvent
		rErr   string
	}{
		{
			"positive case",
			Vault{},
			tEvent{
				EventModel: &EventModel{1, ""},
				applyFn: func(Entry) error {
					return nil
				},
			},
			"",
		},
		{
			"invalid id",
			Vault{
				commitedEvents: []Event{tEvent{EventModel: &EventModel{10, ""}}},
			},
			tEvent{EventModel: &EventModel{1, ""}},
			"invalid event version",
		},
		{
			"err on event apply",
			Vault{
				commitedEvents: []Event{tEvent{EventModel: &EventModel{10, ""}}},
			},
			tEvent{
				EventModel: &EventModel{11, ""},
				applyFn: func(Entry) error {
					return fmt.Errorf("apply err")
				},
			},
			"apply err",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.Vault.Apply(tt.aEvent)
			if err != nil {
				assert.Equal(t, tt.rErr, err.Error())
				return
			}
			assert.Equal(t, "", tt.rErr, "Expecting error: %s", tt.rErr)
			assert.Equal(t, 1, len(tt.Vault.uncomittedEvents))
		})
	}
}

type tEvent struct {
	*EventModel
	applyFn func(Entry) error
}

func (t tEvent) Apply(e Entry) error {
	if t.applyFn == nil {
		return nil
	}

	return t.applyFn(e)
}
