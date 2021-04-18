package sqlx

import (
	"testing"
	"time"

	"github.com/reekoheek/brankas/pkg/vault"
	"github.com/reekoheek/brankas/pkg/vault/entries/acl"
	"github.com/reekoheek/brankas/pkg/vault/entries/login"
	"gotest.tools/assert"
)

func TestMapper_ToEvent(t *testing.T) {
	now := time.Now()

	table := []struct {
		name   string
		aModel mEvent
		xEvent vault.Event
		xErr   string
	}{
		{
			"acl put",
			mEvent{
				Version: 10,
				ID:      "foo",
				At:      mTime(now),
				Kind:    KIND_ACL_PUT,
				Data: mEventData{
					Username: "foo",
					Mode:     acl.READWRITE,
				},
			},
			acl.Put{
				EventModel: vault.NewEventModel(10, "foo", now),
				Username:   "foo",
				Mode:       acl.READWRITE,
			},
			"",
		},
		{
			"login added",
			mEvent{
				Kind: KIND_LOGIN_ADDED,
			},
			login.Added{},
			"",
		},
		{
			"login name changed",
			mEvent{
				Kind: KIND_LOGIN_NAME_CHANGED,
				Data: mEventData{
					Name: "foo",
				},
			},
			login.NameChanged{
				Name: "foo",
			},
			"",
		},
		{
			"login notes changed",
			mEvent{
				Kind: KIND_LOGIN_NOTES_CHANGED,
				Data: mEventData{
					Notes: "foo",
				},
			},
			login.NotesChanged{
				Notes: "foo",
			},
			"",
		},
		{
			"login password changed",
			mEvent{
				Kind: KIND_LOGIN_PASSWORD_CHANGED,
				Data: mEventData{
					Password: "foo",
					Expiry:   now,
				},
			},
			login.PasswordChanged{
				Password: "foo",
				Expiry:   now,
			},
			"",
		},
		{
			"login url added",
			mEvent{
				Kind: KIND_LOGIN_URL_ADDED,
				Data: mEventData{
					URL: "foo",
				},
			},
			login.URLAdded{
				URL: "foo",
			},
			"",
		},
		{
			"login url removed",
			mEvent{
				Kind: KIND_LOGIN_URL_REMOVED,
				Data: mEventData{
					URL: "foo",
				},
			},
			login.URLRemoved{
				URL: "foo",
			},
			"",
		},
		{
			"login username changed",
			mEvent{
				Kind: KIND_LOGIN_USERNAME_CHANGED,
				Data: mEventData{
					Username: "foo",
				},
			},
			login.UsernameChanged{
				Username: "foo",
			},
			"",
		},
		{
			"unknown event",
			mEvent{},
			nil,
			"unknown event",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapper{}

			ev, err := m.toEvent(tt.aModel)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)
			assert.DeepEqual(t, tt.xEvent, ev)
		})
	}

}

func TestMapper_ToMEvent(t *testing.T) {
	now := time.Now()

	table := []struct {
		name   string
		aEvent vault.Event
		xModel mEvent
		xErr   string
	}{
		{
			"acl put",
			acl.Put{
				EventModel: vault.NewEventModel(99, "identity", now),
				Username:   "foo",
				Mode:       acl.READONLY,
			},
			mEvent{
				Version: 99,
				ID:      "identity",
				At:      mTime(now),
				Kind:    KIND_ACL_PUT,
				Data: mEventData{
					Username: "foo",
					Mode:     acl.READONLY,
				},
			},
			"",
		},
		{
			"login added",
			login.Added{},
			mEvent{
				Kind: KIND_LOGIN_ADDED,
				Data: mEventData{},
			},
			"",
		},
		{
			"login name changed",
			login.NameChanged{
				Name: "foo",
			},
			mEvent{
				Kind: KIND_LOGIN_NAME_CHANGED,
				Data: mEventData{
					Name: "foo",
				},
			},
			"",
		},
		{
			"login notes changed",
			login.NotesChanged{
				Notes: "foo",
			},
			mEvent{
				Kind: KIND_LOGIN_NOTES_CHANGED,
				Data: mEventData{
					Notes: "foo",
				},
			},
			"",
		},
		{
			"login password changed",
			login.PasswordChanged{
				Password: "foo",
				Expiry:   now,
			},
			mEvent{
				Kind: KIND_LOGIN_PASSWORD_CHANGED,
				Data: mEventData{
					Password: "foo",
					Expiry:   now,
				},
			},
			"",
		},
		{
			"login url added",
			login.URLAdded{
				URL: "foo",
			},
			mEvent{
				Kind: KIND_LOGIN_URL_ADDED,
				Data: mEventData{
					URL: "foo",
				},
			},
			"",
		},
		{
			"login url removed",
			login.URLRemoved{
				URL: "foo",
			},
			mEvent{
				Kind: KIND_LOGIN_URL_REMOVED,
				Data: mEventData{
					URL: "foo",
				},
			},
			"",
		},
		{
			"login username changed",
			login.UsernameChanged{
				Username: "foo",
			},
			mEvent{
				Kind: KIND_LOGIN_USERNAME_CHANGED,
				Data: mEventData{
					Username: "foo",
				},
			},
			"",
		},
		{
			"unknown event",
			tEvent{},
			mEvent{},
			"unknown event",
		},
		{
			"unknown event",
			nil,
			mEvent{},
			"unknown event",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapper{}

			mo, err := m.toMEvent(tt.aEvent)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)
			assert.Equal(t, tt.xModel, mo)
		})
	}
}
