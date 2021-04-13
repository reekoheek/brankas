package sqlx

import (
	"database/sql/driver"
	"time"

	"github.com/reekoheek/brankas/pkg/vault"
)

type mTime time.Time

func (t *mTime) Scan(v interface{}) error {
	vt, err := time.Parse(time.RFC3339, v.(string))
	if err != nil {
		return err
	}
	*t = mTime(vt)
	return nil
}

func (t mTime) Value() (driver.Value, error) {
	return time.Time(t).Format(time.RFC3339), nil
}

type mEvent struct {
	VaultID string `db:"vault_id"`
	Version int    `db:"version"`
	ID      string `db:"id"`
	At      mTime  `db:"at"`
	Kind    string `db:"kind"`
	Data    string `db:"data"`
}

type toEventMapper interface {
	toEvent(mEvent) (vault.Event, error)
}

type toMEventMapper interface {
	toMEvent(vault.Event) (mEvent, error)
}

type mapper struct {
}

func (m *mapper) toEvent(mEvent) (vault.Event, error) {
	return nil, nil
}

func (m *mapper) toMEvent(vault.Event) (mEvent, error) {
	return mEvent{}, nil
}
