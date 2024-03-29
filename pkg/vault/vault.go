package vault

import "fmt"

type Vault struct {
	id               string
	commitedEvents   []Event
	uncommitedEvents []Event
	entries          map[string]Entry
}

func New(id string, evts []Event) Vault {
	return Vault{
		id:             id,
		entries:        map[string]Entry{},
		commitedEvents: evts,
	}
}

func (v *Vault) ID() string {
	return v.id
}

func (v *Vault) Apply(evt Event) error {
	if v.Version()+1 != evt.Version() {
		return fmt.Errorf("invalid event version")
	}

	id := evt.ID()

	entry := v.entries[id]

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

func (v *Vault) Events(version int) []Event {
	var evs []Event
	for _, ev := range v.commitedEvents {
		if ev.Version() >= version {
			evs = append(evs, ev)
		}
	}
	return evs
}

func (v *Vault) UncommitedEvents() []Event {
	return v.uncommitedEvents
}
