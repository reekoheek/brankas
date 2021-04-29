package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gotest.tools/assert"
)

func TestUser_ChangePassword(t *testing.T) {
	table := []struct {
		name      string
		aPassword string
	}{
		{
			"case 1",
			"foo",
		},
		{
			"case 2",
			"bar",
		},
		{
			"case 3",
			"baz",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			u := User{}

			if err := u.ChangePassword(tt.aPassword); err != nil {
				t.Error(err)
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(tt.aPassword)); err != nil {
				t.Error(err)
				return
			}
		})
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	table := []struct {
		name      string
		dPassword string
		aPassword string
		xErr      string
	}{
		{
			"positive case",
			"foo",
			"foo",
			"",
		},
		{
			"invalid password",
			"foo",
			"bar",
			"invalid password",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := bcrypt.GenerateFromPassword([]byte(tt.dPassword), BCRYPT_COST)
			if err != nil {
				t.Error(err)
				return
			}

			u := User{Password: string(bytes)}

			if err := u.ValidatePassword(tt.aPassword); err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr)
		})
	}
}
