package auth

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	BCRYPT_COST = 14
)

type User struct {
	Username string
	Password string
}

func (u *User) ChangePassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_COST)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) ValidatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}

type RepoUserGetter interface {
	Get(ctx context.Context, username string) (User, error)
}

type UserRepository interface {
	RepoUserGetter
}
