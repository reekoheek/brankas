package auth

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

func TestService_Login(t *testing.T) {
	table := []struct {
		name               string
		dUserGetter        RepoUserGetter
		dPasswordValidator passwordValidator
		aDTO               LoginDTO
		xErr               string
	}{
		{
			"positive case",
			tUserGetter(func(context.Context, string) (User, error) {
				return User{
					Username: "foo",
					Password: "foo",
				}, nil
			}),
			tPasswordValidator(func(string, string) error {
				return nil
			}),
			LoginDTO{
				Username: "foo",
				Password: "foo",
			},
			"",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				rUserGetter:       tt.dUserGetter,
				passwordValidator: tt.dPasswordValidator,
			}

			result, err := s.Login(context.TODO(), tt.aDTO)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr)
			assert.Equal(t, result.AccessToken != "", true)
		})
	}
}

type tUserGetter func(context.Context, string) (User, error)

func (t tUserGetter) Get(ctx context.Context, username string) (User, error) {
	return t(ctx, username)
}

type tPasswordValidator func(string, string) error

func (t tPasswordValidator) Validate(hash string, password string) error {
	return t(hash, password)
}
