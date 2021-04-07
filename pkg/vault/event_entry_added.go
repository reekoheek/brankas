package vault

type EntryAdded struct {
	EventModel
}

func (ev EntryAdded) Apply(e Entry) error {
	// entry := v.Has(ev.index)
	return nil
}
