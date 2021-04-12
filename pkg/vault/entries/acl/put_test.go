package acl

import (
	"testing"

	"github.com/reekoheek/brankas/pkg/vault"
	"gotest.tools/assert"
)

func TestPut_Mutate(t *testing.T) {
	table := []struct {
		name   string
		ev     Put
		aEntry vault.Entry
		xErr   string
	}{
		{
			"positive 1",
			Put{
				Username: "foo",
				Mode:     READWRITE,
			},
			nil,
			"",
		},
		{
			"positive 2",
			Put{
				Username: "foo",
				Mode:     READONLY,
			},
			ACL{
				modes: map[string]int{},
			},
			"",
		},
		{
			"positive 3",
			Put{
				Username: "foo",
				Mode:     NONE,
			},
			nil,
			"",
		},
		{
			"if acl put err",
			Put{
				Username: "foo",
				Mode:     100,
			},
			nil,
			"invalid mode",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := tt.ev.Mutate(tt.aEntry)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)

			e := entry.(ACL)
			assert.Equal(t, tt.ev.Mode, e.modes[tt.ev.Username])
		})
	}
}
