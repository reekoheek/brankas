package login

import (
	"testing"

	"github.com/reekoheek/brankas/pkg/vault"
	"gotest.tools/assert"
)

func TestURLAdded_Mutate(t *testing.T) {
	table := []struct {
		name   string
		ev     URLAdded
		aEntry vault.Entry
		xErr   string
	}{
		{
			"positive case",
			URLAdded{
				URL: "foo",
			},
			Login{},
			"",
		},
		{
			"if url exists",
			URLAdded{
				URL: "foo",
			},
			Login{
				URLs: []string{"foo"},
			},
			"url exists",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			e, err := tt.ev.Mutate(tt.aEntry)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr)
			entry := e.(Login)
			assert.Equal(t, tt.ev.ID(), entry.ID())
			assert.Equal(t, tt.ev.Version(), entry.Version())
			assert.Equal(t, tt.ev.URL, entry.URLs[0])
		})
	}
}
