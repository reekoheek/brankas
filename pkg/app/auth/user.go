package auth

import (
	"context"
)

type User struct {
	Username string
	Password string
}

type RepoUserGetter interface {
	Get(ctx context.Context, username string) (User, error)
}

type UserRepository interface {
	RepoUserGetter
}
