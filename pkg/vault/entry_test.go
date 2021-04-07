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
			name: "positive 1",
			dID:  "1",
		},
		{
			name: "positive 2",
			dID:  "2",
		},
		{
			name: "positive 3",
			dID:  "3",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			m := EntryModel{id: tt.dID}
			assert.Equal(t, tt.dID, m.ID())
		})
	}
}
