package login

import "github.com/reekoheek/brankas/pkg/vault"

type URLAdded struct {
	vault.EventModel
	URL string
}

func (ev URLAdded) Mutate(entry vault.Entry) (vault.Entry, error) {
	e, ok := entry.(Login)
	if !ok {
		return nil, vault.ErrInvalidEntry
	}

	e.EntryModel = vault.NewEntryModel(ev.ID(), ev.Version())
	if err := e.AddURL(ev.URL); err != nil {
		return nil, err
	}

	return e, nil
}
