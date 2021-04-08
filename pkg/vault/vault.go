package vault

import "fmt"

type Vault struct {
	commitedEvents   []Event
	uncommitedEvents []Event
	entries          map[string]Entry
}

func New(evts []Event) Vault {
	return Vault{
		commitedEvents: evts,
	}
}

func (v *Vault) Apply(evt Event) error {
	if v.Version()+1 != evt.Version() {
		return fmt.Errorf("invalid event version")
	}

	id := evt.ID()

	entry := v.Get(id)

	var err error
	if entry, err = evt.Mutate(entry); err != nil {
		return err
	}

	v.entries[id] = entry

	v.uncommitedEvents = append(v.uncommitedEvents, evt)

	return nil
}

func (v *Vault) Version() int {
	if len(v.uncommitedEvents) != 0 {
		return v.uncommitedEvents[len(v.uncommitedEvents)-1].Version()
	}

	if len(v.commitedEvents) != 0 {
		return v.commitedEvents[len(v.commitedEvents)-1].Version()
	}

	return 0
}

func (v *Vault) Get(id string) Entry {
	return v.entries[id]
}
