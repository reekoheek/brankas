package auth

import (
	"testing"

	"gotest.tools/assert"
)

func TestBcryptHashValidator(t *testing.T) {
	table := []struct {
		name string
		p1   string
		p2   string
		xErr string
	}{
		{
			"positive case 1",
			"foo",
			"foo",
			"",
		},
		{
			"positive case 2",
			"bar",
			"bar",
			"",
		},
		{
			"password not same",
			"foo",
			"bar",
			"invalid password",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			b := bcryptHashValidator{}

			hash, err := b.Hash(tt.p1)
			if err != nil {
				t.Error(err)
				return
			}

			if err := b.Validate(hash, tt.p2); err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr)
		})
	}
}
