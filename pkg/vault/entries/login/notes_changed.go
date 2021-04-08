package login

import "github.com/reekoheek/brankas/pkg/vault"

type NotesChanged struct {
	vault.EventModel
	Notes string
}

func (ev NotesChanged) Mutate(entry vault.Entry) (vault.Entry, error) {
	e, ok := entry.(Login)
	if !ok {
		return nil, vault.ErrInvalidEntry
	}

	e.EntryModel = vault.NewEntryModel(ev.ID(), ev.Version())
	e.Notes = ev.Notes

	return e, nil
}
