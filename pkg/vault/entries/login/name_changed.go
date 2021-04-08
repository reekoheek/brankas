package login

import "github.com/reekoheek/brankas/pkg/vault"

type NameChanged struct {
	vault.EventModel
	Name string
}

func (ev NameChanged) Mutate(entry vault.Entry) (vault.Entry, error) {
	e, ok := entry.(Login)
	if !ok {
		return nil, vault.ErrInvalidEntry
	}

	e.EntryModel = vault.NewEntryModel(ev.ID(), ev.Version())
	e.Name = ev.Name

	return e, nil
}
