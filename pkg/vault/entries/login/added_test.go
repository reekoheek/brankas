package login

import (
	"testing"

	"github.com/reekoheek/brankas/pkg/vault"
	"github.com/stretchr/testify/assert"
)

func TestAdded_Mutate(t *testing.T) {
	table := []struct {
		name   string
		ev     vault.Event
		aEntry vault.Entry
		rErr   string
	}{
		{
			"positive case 1",
			Added{EventModel: vault.NewEventModel(99, "foo")},
			nil,
			"",
		},
		{
			"entry exists",
			Added{EventModel: vault.NewEventModel(99, "foo")},
			Login{},
			"entry exists",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			e, err := tt.ev.Mutate(tt.aEntry)
			if err != nil {
				assert.Equal(t, tt.rErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.rErr, "Expecting error: %s", tt.rErr)

			entry := e.(Login)
			assert.Equal(t, "", entry.Name)
		})
	}
}
