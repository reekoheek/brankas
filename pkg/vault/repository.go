package vault

import "context"

type RepoGetter interface {
	Get(ctx context.Context, id string) (Vault, error)
}

type RepoPersister interface {
	Persist(ctx context.Context, v Vault) error
}

type Repository interface {
	RepoGetter
	RepoPersister
}
