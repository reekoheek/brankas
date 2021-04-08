package login

import (
	"fmt"
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

func (e *Login) Expired() bool {
	now := time.Now()
	return now.Equal(e.Expiry) || now.After(e.Expiry)
}

func (e *Login) AddURL(url string) error {
	if e.hasURL(url) {
		return fmt.Errorf("url exists")
	}

	e.URLs = append(e.URLs, url)

	return nil
}

func (e *Login) RemoveURL(url string) error {
	if !e.hasURL(url) {
		return fmt.Errorf("url not exists")
	}

	urls := []string{}
	for _, u := range e.URLs {
		if u != url {
			urls = append(urls, u)
		}
	}

	e.URLs = urls

	return nil
}

func (e *Login) hasURL(url string) bool {
	for _, u := range e.URLs {
		if u == url {
			return true
		}
	}

	return false
}
