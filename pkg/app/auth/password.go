package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	BCRYPT_COST = 14
)

type passwordHasher interface {
	Hash(password string) (string, error)
}

type passwordValidator interface {
	Validate(hash string, password string) error
}

type passwordHashValidator interface {
	passwordHasher
	passwordValidator
}

type bcryptHashValidator struct {
	cost int
}

func (b bcryptHashValidator) Hash(password string) (string, error) {
	bb, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(bb), err
}

func (b bcryptHashValidator) Validate(hash string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}
