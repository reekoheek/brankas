package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/reekoheek/brankas/pkg/vault"
)

const (
	VaultSchema = `
	CREATE TABLE IF NOT EXISTS event (
		vault_id VARCHAR(100),
		version  INT,
		id       VARCHAR(100),
		at       VARCHAR(100),
		kind     VARCHAR(100),
		data     TEXT,
		PRIMARY KEY (vault_id, version)
	);
	`
)

type vaultRepository struct {
	base
	mToEvent  toEventMapper
	mToMEvent toMEventMapper
}

func (r *vaultRepository) Get(ctx context.Context, id string) (vault.Vault, error) {
	tx, err := r.tx(ctx)
	if err != nil {
		return vault.Vault{}, err
	}

	mevs := []mEvent{}
	if err := tx.SelectContext(ctx, &mevs, `SELECT * FROM event WHERE vault_id = ?`, id); err != nil {
		return vault.Vault{}, err
	}

	evs := []vault.Event{}
	for _, mev := range mevs {
		ev, err := r.mToEvent.toEvent(mev)
		if err != nil {
			return vault.Vault{}, err
		}
		evs = append(evs, ev)
	}

	return vault.New(id, evs), nil
}

func (r *vaultRepository) Persist(ctx context.Context, v vault.Vault) error {
	tx, err := r.tx(ctx)
	if err != nil {
		return err
	}

	sql := `INSERT INTO event
			(vault_id, version, id, at, kind, data)
		VALUES
			(:vault_id, :version, :id, :at, :kind, :data)`

	evs := v.UncommitedEvents()
	for _, ev := range evs {
		mev, err := r.mToMEvent.toMEvent(ev)
		if err != nil {
			return err
		}
		mev.VaultID = v.ID()

		if _, err = tx.NamedExecContext(ctx, sql, mev); err != nil {
			return err
		}
	}

	return nil
}

func NewVaultRepository(db *sqlx.DB) vault.Repository {
	m := &mapper{}
	return &vaultRepository{
		base:      base{db},
		mToEvent:  m,
		mToMEvent: m,
	}
}
