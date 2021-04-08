package acl

import (
	"fmt"

	"github.com/reekoheek/brankas/pkg/vault"
)

const (
	NONE      = 0
	READONLY  = 1
	READWRITE = 2
)

type ACL struct {
	vault.EntryModel
	modes map[string]int
}

func (acl ACL) WriteBy(username string) bool {
	mode := acl.modes[username]
	return mode == READWRITE
}

func (acl ACL) ReadBy(username string) bool {
	mode := acl.modes[username]
	return mode == READWRITE || mode == READONLY
}

func (acl ACL) Put(username string, mode int) error {
	if mode < NONE || mode > READWRITE {
		return fmt.Errorf("invalid mode")
	}

	if mode == NONE {
		delete(acl.modes, username)
	} else {
		acl.modes[username] = mode
	}

	return nil
}
