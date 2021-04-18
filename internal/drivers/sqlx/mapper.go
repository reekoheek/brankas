package sqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/reekoheek/brankas/pkg/vault"
	"github.com/reekoheek/brankas/pkg/vault/entries/acl"
	"github.com/reekoheek/brankas/pkg/vault/entries/login"
)

const (
	KIND_ACL_PUT                = "ACLPut"
	KIND_LOGIN_ADDED            = "LoginAdded"
	KIND_LOGIN_NAME_CHANGED     = "LoginNameChanged"
	KIND_LOGIN_NOTES_CHANGED    = "LoginNotesChanged"
	KIND_LOGIN_PASSWORD_CHANGED = "LoginPasswordChanged"
	KIND_LOGIN_URL_ADDED        = "LoginURLAdded"
	KIND_LOGIN_URL_REMOVED      = "LoginURLRemoved"
	KIND_LOGIN_USERNAME_CHANGED = "LoginUsernameChanged"
)

type mTime time.Time

func (t *mTime) Scan(v interface{}) error {
	vt, err := time.Parse(time.RFC3339, v.(string))
	if err != nil {
		return err
	}
	*t = mTime(vt)
	return nil
}

func (t mTime) Value() (driver.Value, error) {
	return time.Time(t).Format(time.RFC3339), nil
}

type mEvent struct {
	VaultID string     `db:"vault_id"`
	Version int        `db:"version"`
	ID      string     `db:"id"`
	At      mTime      `db:"at"`
	Kind    string     `db:"kind"`
	Data    mEventData `db:"data"`
}

type mEventData struct {
	Username string
	Mode     int
	Name     string
	Notes    string
	Password string
	Expiry   time.Time
	URL      string
}

func (d *mEventData) Scan(v interface{}) error {
	data := mEventData{}

	if v == nil {
		*d = data
		return nil
	}

	if err := json.Unmarshal([]byte(v.(string)), &data); err != nil {
		return fmt.Errorf("invalid data")
	}

	*d = data
	return nil
}

func (d mEventData) Value() (driver.Value, error) {
	bb, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return string(bb), nil
}

type toEventMapper interface {
	toEvent(mEvent) (vault.Event, error)
}

type toMEventMapper interface {
	toMEvent(vault.Event) (mEvent, error)
}

type mapper struct {
}

func (m *mapper) toEvent(me mEvent) (vault.Event, error) {
	em := vault.NewEventModel(me.Version, me.ID, time.Time(me.At))

	switch me.Kind {
	case KIND_ACL_PUT:
		return acl.Put{
			EventModel: em,
			Username:   me.Data.Username,
			Mode:       me.Data.Mode,
		}, nil
	case KIND_LOGIN_ADDED:
		return login.Added{
			EventModel: em,
		}, nil
	case KIND_LOGIN_NAME_CHANGED:
		return login.NameChanged{
			EventModel: em,
			Name:       me.Data.Name,
		}, nil
	case KIND_LOGIN_NOTES_CHANGED:
		return login.NotesChanged{
			EventModel: em,
			Notes:      me.Data.Notes,
		}, nil
	case KIND_LOGIN_PASSWORD_CHANGED:
		return login.PasswordChanged{
			EventModel: em,
			Password:   me.Data.Password,
			Expiry:     me.Data.Expiry,
		}, nil
	case KIND_LOGIN_URL_ADDED:
		return login.URLAdded{
			EventModel: em,
			URL:        me.Data.URL,
		}, nil
	case KIND_LOGIN_URL_REMOVED:
		return login.URLRemoved{
			EventModel: em,
			URL:        me.Data.URL,
		}, nil
	case KIND_LOGIN_USERNAME_CHANGED:
		return login.UsernameChanged{
			EventModel: em,
			Username:   me.Data.Username,
		}, nil
	}
	return nil, fmt.Errorf("unknown event")
}

func (m *mapper) toMEvent(ev vault.Event) (mEvent, error) {
	if ev == nil {
		return mEvent{}, fmt.Errorf("unknown event")
	}

	mo := mEvent{
		Version: ev.Version(),
		ID:      ev.ID(),
		At:      mTime(ev.At()),
	}

	switch evt := ev.(type) {
	case acl.Put:
		mo.Kind = KIND_ACL_PUT
		mo.Data = mEventData{
			Username: evt.Username,
			Mode:     evt.Mode,
		}
	case login.Added:
		mo.Kind = KIND_LOGIN_ADDED
	case login.NameChanged:
		mo.Kind = KIND_LOGIN_NAME_CHANGED
		mo.Data = mEventData{
			Name: evt.Name,
		}
	case login.NotesChanged:
		mo.Kind = KIND_LOGIN_NOTES_CHANGED
		mo.Data = mEventData{
			Notes: evt.Notes,
		}
	case login.PasswordChanged:
		mo.Kind = KIND_LOGIN_PASSWORD_CHANGED
		mo.Data = mEventData{
			Password: evt.Password,
			Expiry:   evt.Expiry,
		}
	case login.URLAdded:
		mo.Kind = KIND_LOGIN_URL_ADDED
		mo.Data = mEventData{
			URL: evt.URL,
		}
	case login.URLRemoved:
		mo.Kind = KIND_LOGIN_URL_REMOVED
		mo.Data = mEventData{
			URL: evt.URL,
		}
	case login.UsernameChanged:
		mo.Kind = KIND_LOGIN_USERNAME_CHANGED
		mo.Data = mEventData{
			Username: evt.Username,
		}
	}

	if mo.Kind == "" {
		return mo, fmt.Errorf("unknown event")
	}

	return mo, nil
}
