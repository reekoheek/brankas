package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/reekoheek/brankas/internal/tx"
)

type base struct {
	db *sqlx.DB
}

func (b *base) tx(ctx context.Context) (*sqlx.Tx, error) {
	t, err := tx.Acquire(ctx, b.db, func() (tx.CommitRollbacker, error) {
		return b.db.Beginx()
	})

	if err != nil {
		return nil, err
	}

	return t.(*sqlx.Tx), nil
}
