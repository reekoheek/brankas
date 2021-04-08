package acl

import "github.com/reekoheek/brankas/pkg/vault"

type Put struct {
	vault.EventModel
	Username string
	Mode     int
}

func (ev Put) Mutate(entry vault.Entry) (vault.Entry, error) {
	e := ACL{
		EntryModel: vault.NewEntryModel(ev.ID(), ev.Version()),
		modes:      map[string]int{},
	}

	if entry != nil {
		e := entry.(ACL)
		e.EntryModel = vault.NewEntryModel(ev.ID(), ev.Version())
	}

	if err := e.Put(ev.Username, ev.Mode); err != nil {
		return nil, err
	}

	return e, nil
}
