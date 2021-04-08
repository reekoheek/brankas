package login

import (
	"testing"

	"github.com/reekoheek/brankas/pkg/vault"
	"github.com/stretchr/testify/assert"
)

func TestUsernameChanged_Mutate(t *testing.T) {
	table := []struct {
		name   string
		ev     UsernameChanged
		aEntry vault.Entry
		rErr   string
	}{
		{
			"positive case",
			UsernameChanged{
				Username: "foo",
			},
			Login{},
			"",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			e, err := tt.ev.Mutate(tt.aEntry)
			if err != nil {
				assert.Equal(t, tt.rErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.rErr)
			entry := e.(Login)
			assert.Equal(t, tt.ev.ID(), entry.ID())
			assert.Equal(t, tt.ev.Version(), entry.Version())
			assert.Equal(t, tt.ev.Username, entry.Username)
		})
	}
}