package login

import (
	"testing"
	"time"

	"github.com/reekoheek/brankas/pkg/vault"
	"gotest.tools/assert"
)

func TestAdded_Mutate(t *testing.T) {
	table := []struct {
		name   string
		ev     vault.Event
		aEntry vault.Entry
		xErr   string
	}{
		{
			"positive case 1",
			Added{EventModel: vault.NewEventModel(99, "foo", time.Now())},
			nil,
			"",
		},
		{
			"entry exists",
			Added{EventModel: vault.NewEventModel(99, "foo", time.Now())},
			Login{},
			"entry exists",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			e, err := tt.ev.Mutate(tt.aEntry)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "Expecting error: %s", tt.xErr)

			entry := e.(Login)
			assert.Equal(t, "", entry.Name)
		})
	}
}
