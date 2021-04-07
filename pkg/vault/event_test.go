package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventModel_Index(t *testing.T) {
	table := []struct {
		name     string
		dVersion int
	}{
		{
			name:     "positive 1",
			dVersion: 1,
		},
		{
			name:     "positive 2",
			dVersion: 2,
		},
		{
			name:     "positive 3",
			dVersion: 3,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			m := EventModel{
				version: tt.dVersion,
			}

			assert.Equal(t, tt.dVersion, m.Version())
		})
	}
}
