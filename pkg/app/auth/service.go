package auth

import (
	"context"

	"github.com/dgrijalva/jwt-go"
)

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResultDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Loginer interface {
	Login(context.Context, LoginDTO) (LoginResultDTO, error)
}

type Service interface {
	Loginer
}

type service struct {
	jwtSignatureKey   []byte
	rUserGetter       RepoUserGetter
	passwordValidator passwordValidator
}

func New(jwtSignatureKey []byte) Service {
	hv := bcryptHashValidator{}
	return &service{
		jwtSignatureKey:   jwtSignatureKey,
		passwordValidator: hv,
	}
}

func (s *service) Login(ctx context.Context, dto LoginDTO) (LoginResultDTO, error) {
	user, err := s.rUserGetter.Get(ctx, dto.Username)
	if err != nil {
		return LoginResultDTO{}, err
	}

	if err := s.passwordValidator.Validate(user.Password, dto.Password); err != nil {
		return LoginResultDTO{}, err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return LoginResultDTO{}, err
	}

	return LoginResultDTO{
		AccessToken: accessToken,
	}, nil
}

func (s *service) generateAccessToken(user User) (string, error) {
	claims := jwt.StandardClaims{
		Subject: user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtSignatureKey)
}
