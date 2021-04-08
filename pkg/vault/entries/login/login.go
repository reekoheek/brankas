package login

import (
	"time"

	"github.com/reekoheek/brankas/pkg/vault"
)

type Login struct {
	vault.EntryModel
	Name     string
	Username string
	Password string
	Expiry   time.Time
	URLs     []string
	Notes    string
}

func (e Login) Expired() bool {
	now := time.Now()
	return now.Equal(e.Expiry) || now.After(e.Expiry)
}
