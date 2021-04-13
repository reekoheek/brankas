package sqlx

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/reekoheek/brankas/internal/tx"
	"gotest.tools/assert"
)

func TestBase_Tx(t *testing.T) {
	table := []struct {
		name string
		ctx  context.Context
		xErr string
	}{
		{
			"positive case 1",
			tx.Context(context.TODO(), tx.New()),
			"",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			db := sqlx.MustConnect("sqlite3", ":memory:")
			defer db.Close()

			b := &base{db}

			tx, err := b.tx(tt.ctx)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)
			assert.Check(t, tx != nil)
		})
	}
}
