package sync

import (
	"context"
	"fmt"
	"testing"

	"github.com/reekoheek/brankas/pkg/vault"
	"gotest.tools/assert"
)

func TestService_Push(t *testing.T) {
	table := []struct {
		name          string
		aDTO          PushDTO
		dGetterErr    error
		dPersisterErr error
		dMapperErr    error
		xErr          string
	}{
		{
			"positive",
			PushDTO{
				ID:     "id",
				Events: []EventDTO{{Version: 1}},
			},
			nil,
			nil,
			nil,
			"",
		},
		{
			"if invalid id",
			PushDTO{},
			nil,
			nil,
			nil,
			"invalid id",
		},
		{
			"if invalid events",
			PushDTO{
				ID: "id",
			},
			nil,
			nil,
			nil,
			"invalid events",
		},
		{
			"if getter err",
			PushDTO{
				ID:     "id",
				Events: []EventDTO{{Version: 1}},
			},
			fmt.Errorf("getter err"),
			nil,
			nil,
			"getter err",
		},
		{
			"if mapper err",
			PushDTO{
				ID:     "id",
				Events: []EventDTO{{Version: 1}},
			},
			nil,
			fmt.Errorf("mapper err"),
			nil,
			"mapper err",
		},
		{
			"if persister err",
			PushDTO{
				ID:     "id",
				Events: []EventDTO{{Version: 1}},
			},
			nil,
			nil,
			fmt.Errorf("persister err"),
			"persister err",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				rVaultGetter: tGetter(func(context.Context, string) (vault.Vault, error) {
					if tt.dGetterErr != nil {
						return vault.Vault{}, tt.dGetterErr
					}
					return vault.New(nil), nil
				}),
				rVaultPersister: tPersister(func(ctx context.Context, v vault.Vault) error {
					if tt.dPersisterErr != nil {
						return tt.dPersisterErr
					}
					evs := v.UncommitedEvents()
					assert.Equal(t, len(tt.aDTO.Events), len(evs))
					return nil
				}),
				mToEvent: tToEventMapper(func(dto EventDTO) (vault.Event, error) {
					if tt.dMapperErr != nil {
						return nil, tt.dMapperErr
					}
					return tEvent{
						EventModel: vault.NewEventModel(dto.Version, dto.ID, dto.At),
					}, nil
				}),
			}

			err := s.Push(context.TODO(), tt.aDTO)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)
		})
	}
}

type tGetter func(ctx context.Context, id string) (vault.Vault, error)

func (t tGetter) Get(ctx context.Context, id string) (vault.Vault, error) {
	return t(ctx, id)
}

type tPersister func(context.Context, vault.Vault) error

func (t tPersister) Persist(ctx context.Context, v vault.Vault) error {
	return t(ctx, v)
}

type tToEventMapper func(dto EventDTO) (vault.Event, error)

func (t tToEventMapper) ToEvent(dto EventDTO) (vault.Event, error) {
	return t(dto)
}
