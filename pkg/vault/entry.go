package vault

type Entry interface {
	ID() string
}

type EntryModel struct {
	id string
}

func (m *EntryModel) ID() string {
	return m.id
}
