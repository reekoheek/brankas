package auth

import "context"

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenDTO struct {
	Token string `json:"token"`
}

type Loginer interface {
	Login(context.Context, LoginDTO) (TokenDTO, error)
}

type Service interface {
	Loginer
}

type service struct {
}

func New() Service {
	return &service{}
}

func (s *service) Login(ctx context.Context, dto LoginDTO) (TokenDTO, error) {
	return TokenDTO{}, nil
}
