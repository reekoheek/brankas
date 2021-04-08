package vault

import "fmt"

var ErrEntryExists = fmt.Errorf("entry exists")
var ErrInvalidEntry = fmt.Errorf("invalid entry")

type Entry interface {
	ID() string
	Version() int
}

type EntryModel struct {
	id      string
	version int
}

func (m EntryModel) ID() string {
	return m.id
}

func (m EntryModel) Version() int {
	return m.version
}

func NewEntryModel(id string, version int) EntryModel {
	return EntryModel{id, version}
}
