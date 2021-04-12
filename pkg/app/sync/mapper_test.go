package sync

import (
	"testing"
	"time"

	"github.com/reekoheek/brankas/pkg/vault"
	"github.com/reekoheek/brankas/pkg/vault/entries/acl"
	"github.com/reekoheek/brankas/pkg/vault/entries/login"
	"gotest.tools/assert"
)

func TestMapper_ToDTO(t *testing.T) {
	now := time.Now()
	m := NewMapper()

	table := []struct {
		name   string
		aEvent vault.Event
		xDTO   EventDTO
		xErr   string
	}{
		{
			"acl put",
			acl.Put{
				EventModel: vault.NewEventModel(1, "id", now),
				Username:   "foo",
				Mode:       acl.READONLY,
			},
			EventDTO{
				Kind:     KIND_ACL_PUT,
				ID:       "id",
				Version:  1,
				At:       now,
				Username: "foo",
				Mode:     acl.READONLY,
			},
			"",
		},
		{
			"login added",
			login.Added{},
			EventDTO{
				Kind: KIND_LOGIN_ADDED,
			},
			"",
		},
		{
			"login name changed",
			login.NameChanged{
				Name: "foo",
			},
			EventDTO{
				Kind: KIND_LOGIN_NAME_CHANGED,
				Name: "foo",
			},
			"",
		},
		{
			"login notes changed",
			login.NotesChanged{
				Notes: "foo",
			},
			EventDTO{
				Kind:  KIND_LOGIN_NOTES_CHANGED,
				Notes: "foo",
			},
			"",
		},
		{
			"login password changed",
			login.PasswordChanged{
				Password: "foo",
				Expiry:   now,
			},
			EventDTO{
				Kind:     KIND_LOGIN_PASSWORD_CHANGED,
				Password: "foo",
				Expiry:   now,
			},
			"",
		},
		{
			"login url added",
			login.URLAdded{
				URL: "foo",
			},
			EventDTO{
				Kind: KIND_LOGIN_URL_ADDED,
				URL:  "foo",
			},
			"",
		},
		{
			"login url removed",
			login.URLRemoved{
				URL: "foo",
			},
			EventDTO{
				Kind: KIND_LOGIN_URL_REMOVED,
				URL:  "foo",
			},
			"",
		},
		{
			"login username changed",
			login.UsernameChanged{
				Username: "foo",
			},
			EventDTO{
				Kind:     KIND_LOGIN_USERNAME_CHANGED,
				Username: "foo",
			},
			"",
		},
		{
			"unknown event",
			tEvent{},
			EventDTO{},
			"unknown event",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			dto, err := m.ToDTO(tt.aEvent)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)
			assert.DeepEqual(t, tt.xDTO, dto)
		})
	}
}

func TestMapper_ToEvent(t *testing.T) {
	now := time.Now()
	m := NewMapper()

	table := []struct {
		name   string
		aDTO   EventDTO
		xEvent vault.Event
		xErr   string
	}{
		{
			"acl put",
			EventDTO{
				Kind:     KIND_ACL_PUT,
				ID:       "id",
				Version:  1,
				At:       now,
				Username: "foo",
				Mode:     acl.READONLY,
			},
			acl.Put{
				EventModel: vault.NewEventModel(1, "id", now),
				Username:   "foo",
				Mode:       acl.READONLY,
			},
			"",
		},
		{
			"login name changed",
			EventDTO{
				Kind: KIND_LOGIN_NAME_CHANGED,
				Name: "foo",
			},
			login.NameChanged{
				Name: "foo",
			},
			"",
		},
		{
			"login notes changed",
			EventDTO{
				Kind:  KIND_LOGIN_NOTES_CHANGED,
				Notes: "foo",
			},
			login.NotesChanged{
				Notes: "foo",
			},
			"",
		},
		{
			"login password changed",
			EventDTO{
				Kind:     KIND_LOGIN_PASSWORD_CHANGED,
				Password: "foo",
			},
			login.PasswordChanged{
				Password: "foo",
			},
			"",
		},
		{
			"login url added",
			EventDTO{
				Kind: KIND_LOGIN_URL_ADDED,
				URL:  "foo",
			},
			login.URLAdded{
				URL: "foo",
			},
			"",
		},
		{
			"login url removed",
			EventDTO{
				Kind: KIND_LOGIN_URL_REMOVED,
				URL:  "foo",
			},
			login.URLRemoved{
				URL: "foo",
			},
			"",
		},
		{
			"login username changed",
			EventDTO{
				Kind:     KIND_LOGIN_USERNAME_CHANGED,
				Username: "foo",
			},
			login.UsernameChanged{
				Username: "foo",
			},
			"",
		},
		{
			"unknown event",
			EventDTO{
				Kind: "unknown",
			},
			nil,
			"unknown event",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := m.ToEvent(tt.aDTO)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)
			assert.DeepEqual(t, tt.xEvent, ev)
		})
	}
}

type tEvent struct {
	vault.EventModel
}

func (t tEvent) Mutate(vault.Entry) (vault.Entry, error) {
	return nil, nil
}
