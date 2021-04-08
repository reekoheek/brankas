package login

import "github.com/reekoheek/brankas/pkg/vault"

type URLRemoved struct {
	vault.EventModel
	URL string
}

func (ev URLRemoved) Mutate(entry vault.Entry) (vault.Entry, error) {
	e, ok := entry.(Login)
	if !ok {
		return nil, vault.ErrInvalidEntry
	}

	e.EntryModel = vault.NewEntryModel(ev.ID(), ev.Version())
	if err := e.RemoveURL(ev.URL); err != nil {
		return nil, err
	}

	return e, nil
}
