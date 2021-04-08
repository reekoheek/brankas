package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventModel_Version(t *testing.T) {
	table := []struct {
		name     string
		dVersion int
	}{
		{
			"positive 1",
			1,
		},
		{
			"positive 2",
			2,
		},
		{
			"positive 3",
			3,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			m := EventModel{version: tt.dVersion}

			assert.Equal(t, tt.dVersion, m.Version())
		})
	}
}

func TestEventModel_ID(t *testing.T) {
	table := []struct {
		name string
		dID  string
	}{
		{
			"positive 1",
			"1",
		},
		{
			"positive 2",
			"2",
		},
		{
			"positive 3",
			"3",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			m := EventModel{id: tt.dID}
			assert.Equal(t, tt.dID, m.ID())
		})
	}
}

func TestNewEventModel(t *testing.T) {
	table := []struct {
		name     string
		aVersion int
		aID      string
	}{
		{
			"positive case 1",
			1,
			"foo",
		},
		{
			"positive case 2",
			2,
			"foo",
		},
		{
			"positive case 3",
			1,
			"bar",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			m := NewEventModel(tt.aVersion, tt.aID)
			assert.Equal(t, tt.aVersion, m.version)
			assert.Equal(t, tt.aID, m.id)
		})
	}
}
