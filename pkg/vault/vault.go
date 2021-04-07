package vault

import "fmt"

type Vault struct {
	commitedEvents   []Event
	uncomittedEvents []Event
	entries          map[string]Entry
}

// New Vault
func New(evts []Event) Vault {
	return Vault{
		commitedEvents: evts,
		entries:        map[string]Entry{},
	}
}

// Apply method
func (v *Vault) Apply(evt Event) error {
	if v.Version()+1 != evt.Version() {
		return fmt.Errorf("invalid event version")
	}

	entry := v.Get(evt.ID())

	if err := evt.Apply(entry); err != nil {
		return err
	}

	v.uncomittedEvents = append(v.uncomittedEvents, evt)

	return nil
}

// Version method
func (v *Vault) Version() int {
	if len(v.uncomittedEvents) != 0 {
		return v.uncomittedEvents[len(v.uncomittedEvents)-1].Version()
	}

	if len(v.commitedEvents) != 0 {
		return v.commitedEvents[len(v.commitedEvents)-1].Version()
	}

	return 0
}

// Get method
func (v *Vault) Get(id string) Entry {
	return nil
	// entry := v.entries[id]
	// if entry != nil {
	// 	return entry, nil
	// }

	// for _, ev := range v.commitedEvents {

	// }
}
