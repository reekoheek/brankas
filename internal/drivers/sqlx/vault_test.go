package sqlx

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/reekoheek/brankas/internal/tx"
	"github.com/reekoheek/brankas/pkg/vault"
	"gotest.tools/assert"
)

func TestVaultRepository_Get(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer db.Close()
	db.MustExec(VaultSchema)

	sql := "INSERT INTO event(vault_id, version, id, at, kind, data) VALUES(:vault_id, :version, :id, :at, :kind, :data)"
	db.NamedExec(sql, mEvent{
		VaultID: "bar",
		Version: 8,
		At:      mTime(time.Now()),
		Kind:    "Noop",
	})
	db.NamedExec(sql, mEvent{
		VaultID: "bar",
		Version: 9,
		At:      mTime(time.Now()),
		Kind:    "Noop",
	})
	db.NamedExec(sql, mEvent{
		VaultID: "bar",
		Version: 10,
		At:      mTime(time.Now()),
		Kind:    "Noop",
	})

	table := []struct {
		name     string
		aID      string
		xVersion int
		xLen     int
		xErr     string
	}{
		{
			"case 1",
			"foo",
			0,
			0,
			"",
		},
		{
			"case 2",
			"bar",
			10,
			3,
			"",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := &vaultRepository{
				base: base{db},
				mToEvent: tToEventMapper(func(m mEvent) (vault.Event, error) {
					return tEvent{EventModel: vault.NewEventModel(m.Version, m.ID, time.Time(m.At))}, nil
				}),
			}

			tx.Run(context.TODO(), func(ctx context.Context) error {
				v, err := r.Get(ctx, tt.aID)
				if err != nil {
					assert.Equal(t, tt.xErr, err.Error())
					return nil
				}

				assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)
				assert.Equal(t, tt.xVersion, v.Version())
				assert.Equal(t, tt.xLen, len(v.Events(0)))

				return nil
			})
		})
	}
}

func TestVaultRepository_Persist(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer db.Close()
	db.MustExec(VaultSchema)

	table := []struct {
		name string
		dID  string
		dEvs []vault.Event
		xLen int
		xErr string
	}{
		{
			"positive case",
			"foo",
			[]vault.Event{
				tEvent{EventModel: vault.NewEventModel(1, "", time.Now())},
				tEvent{EventModel: vault.NewEventModel(2, "", time.Now())},
				tEvent{EventModel: vault.NewEventModel(3, "", time.Now())},
			},
			3,
			"",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := &vaultRepository{
				base: base{db},
				mToMEvent: tToMEventMapper(func(ev vault.Event) (mEvent, error) {
					return mEvent{
						Version: ev.Version(),
						ID:      ev.ID(),
						At:      mTime(ev.At()),
					}, nil
				}),
			}

			tx.Run(context.TODO(), func(ctx context.Context) error {
				v := vault.New(tt.dID, nil)
				for _, ev := range tt.dEvs {
					if err := v.Apply(ev); err != nil {
						panic(err)
					}
				}

				if err := r.Persist(ctx, v); err != nil {
					assert.Equal(t, tt.xErr, err.Error())
					return err
				}

				assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)

				return nil
			})

			rows := []mEvent{}
			if err := db.Select(&rows, "SELECT * FROM event WHERE vault_id = ?", tt.dID); err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, tt.xLen, len(rows))
		})
	}
}

type tToEventMapper func(mEvent) (vault.Event, error)

func (t tToEventMapper) toEvent(m mEvent) (vault.Event, error) {
	return t(m)
}

type tToMEventMapper func(vault.Event) (mEvent, error)

func (t tToMEventMapper) toMEvent(ev vault.Event) (mEvent, error) {
	return t(ev)
}

type tEvent struct {
	vault.EventModel
}

func (t tEvent) Mutate(e vault.Entry) (vault.Entry, error) {
	return nil, nil
}
