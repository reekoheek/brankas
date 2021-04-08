package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryModel_ID(t *testing.T) {
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
			m := EntryModel{id: tt.dID}
			assert.Equal(t, tt.dID, m.ID())
		})
	}
}

func TestEntryModel_Version(t *testing.T) {
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
			m := EntryModel{version: tt.dVersion}
			assert.Equal(t, tt.dVersion, m.Version())
		})
	}
}

func TestNewEntryModel(t *testing.T) {
	table := []struct {
		name     string
		aID      string
		aVersion int
	}{
		{
			"positive 1",
			"foo",
			1,
		},
		{
			"positive 2",
			"foo",
			5,
		},
		{
			"positive 2",
			"bar",
			9,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			m := NewEntryModel(tt.aID, tt.aVersion)
			assert.Equal(t, tt.aID, m.id)
			assert.Equal(t, tt.aVersion, m.version)
		})
	}
}
