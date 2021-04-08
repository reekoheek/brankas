package login

import (
	"time"

	"github.com/reekoheek/brankas/pkg/vault"
)

type PasswordChanged struct {
	vault.EventModel
	Password string
	Expiry   time.Time
}

func (ev PasswordChanged) Mutate(entry vault.Entry) (vault.Entry, error) {
	e, ok := entry.(Login)
	if !ok {
		return nil, vault.ErrInvalidEntry
	}

	e.EntryModel = vault.NewEntryModel(ev.ID(), ev.Version())
	e.Password = ev.Password
	e.Expiry = ev.Expiry

	return e, nil
}
