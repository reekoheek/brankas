package login

import "github.com/reekoheek/brankas/pkg/vault"

type UsernameChanged struct {
	vault.EventModel
	Username string
}

func (ev UsernameChanged) Mutate(entry vault.Entry) (vault.Entry, error) {
	e, ok := entry.(Login)
	if !ok {
		return nil, vault.ErrInvalidEntry
	}

	e.EntryModel = vault.NewEntryModel(ev.ID(), ev.Version())
	e.Username = ev.Username

	return e, nil
}
