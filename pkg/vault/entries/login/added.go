package login

import (
	"github.com/reekoheek/brankas/pkg/vault"
)

type Added struct {
	vault.EventModel
}

func (ev Added) Mutate(entry vault.Entry) (vault.Entry, error) {
	if entry != nil {
		return nil, vault.ErrEntryExists
	}

	return Login{
		EntryModel: vault.NewEntryModel(ev.ID(), ev.Version()),
	}, nil
}
