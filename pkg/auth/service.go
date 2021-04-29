package auth

import (
	"context"
	"fmt"
)

type Loginer interface {
	Login(ctx context.Context, username string, password string) (Token, error)
}

type Authenticator interface {
	Authenticate(ctx context.Context, token Token) (User, error)
}

type Service interface {
	Loginer
	Authenticator
}

func NewService() Service {
	return &service{}
}

type service struct {
	rUserGetter RepoUserGetter
}

func (s *service) Login(ctx context.Context, username string, password string) (Token, error) {
	user, err := s.rUserGetter.Get(ctx, username)
	if err != nil {
		return Token{}, err
	}

	return s.generateToken(user), nil
}

func (s *service) Authenticate(ctx context.Context, token Token) (User, error) {
	return User{}, fmt.Errorf("Unimplemented")
}

func (s *service) generateToken(user User) Token {
	return Token{}
}
