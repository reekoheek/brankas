package sync

import (
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

type EventDTO struct {
	Kind     string
	ID       string
	Version  int
	At       time.Time
	Username string
	Mode     int
	Name     string
	Notes    string
	Password string
	Expiry   time.Time
	URL      string
}

type ToDTOMapper interface {
	ToDTO(vault.Event) (EventDTO, error)
}

type ToEventMapper interface {
	ToEvent(EventDTO) (vault.Event, error)
}

type Mapper interface {
	ToDTOMapper
	ToEventMapper
}

type mapper struct{}

func (m *mapper) ToDTO(v vault.Event) (EventDTO, error) {
	dto := EventDTO{
		ID:      v.ID(),
		Version: v.Version(),
		At:      v.At(),
	}

	switch vt := v.(type) {
	case acl.Put:
		dto.Kind = KIND_ACL_PUT
		dto.Username = vt.Username
		dto.Mode = vt.Mode
	case login.Added:
		dto.Kind = KIND_LOGIN_ADDED
	case login.NameChanged:
		dto.Kind = KIND_LOGIN_NAME_CHANGED
		dto.Name = vt.Name
	case login.NotesChanged:
		dto.Kind = KIND_LOGIN_NOTES_CHANGED
		dto.Notes = vt.Notes
	case login.PasswordChanged:
		dto.Kind = KIND_LOGIN_PASSWORD_CHANGED
		dto.Password = vt.Password
		dto.Expiry = vt.Expiry
	case login.URLAdded:
		dto.Kind = KIND_LOGIN_URL_ADDED
		dto.URL = vt.URL
	case login.URLRemoved:
		dto.Kind = KIND_LOGIN_URL_REMOVED
		dto.URL = vt.URL
	case login.UsernameChanged:
		dto.Kind = KIND_LOGIN_USERNAME_CHANGED
		dto.Username = vt.Username
	}

	if dto.Kind == "" {
		return EventDTO{}, fmt.Errorf("unknown event")
	}

	return dto, nil
}

func (m *mapper) ToEvent(dto EventDTO) (vault.Event, error) {
	model := vault.NewEventModel(dto.Version, dto.ID, dto.At)

	switch dto.Kind {
	case KIND_ACL_PUT:
		return acl.Put{
			EventModel: model,
			Username:   dto.Username,
			Mode:       dto.Mode,
		}, nil
	case KIND_LOGIN_NAME_CHANGED:
		return login.NameChanged{
			EventModel: model,
			Name:       dto.Name,
		}, nil
	case KIND_LOGIN_NOTES_CHANGED:
		return login.NotesChanged{
			EventModel: model,
			Notes:      dto.Notes,
		}, nil
	case KIND_LOGIN_PASSWORD_CHANGED:
		return login.PasswordChanged{
			EventModel: model,
			Password:   dto.Password,
			Expiry:     dto.Expiry,
		}, nil
	case KIND_LOGIN_URL_ADDED:
		return login.URLAdded{
			EventModel: model,
			URL:        dto.URL,
		}, nil
	case KIND_LOGIN_URL_REMOVED:
		return login.URLRemoved{
			EventModel: model,
			URL:        dto.URL,
		}, nil
	case KIND_LOGIN_USERNAME_CHANGED:
		return login.UsernameChanged{
			EventModel: model,
			Username:   dto.Username,
		}, nil
	}

	return nil, fmt.Errorf("unknown event")
}

func NewMapper() Mapper {
	return &mapper{}
}
