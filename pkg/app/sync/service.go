package sync

import (
	"context"
	"fmt"

	"github.com/reekoheek/brankas/internal/tx"
	"github.com/reekoheek/brankas/pkg/vault"
)

type PushDTO struct {
	ID     string
	Events []EventDTO
}

type PullDTO struct {
	ID      string
	Version int
}

type service struct {
	rVaultGetter    vault.RepoGetter
	rVaultPersister vault.RepoPersister
	mToDTO          ToDTOMapper
	mToEvent        ToEventMapper
}

func New() *service {
	return &service{}
}

func (s *service) Push(ctx context.Context, dto PushDTO) error {
	if dto.ID == "" {
		return fmt.Errorf("invalid id")
	}

	if len(dto.Events) == 0 {
		return fmt.Errorf("invalid events")
	}

	return tx.Run(ctx, func(ctx context.Context) error {
		v, err := s.rVaultGetter.Get(ctx, dto.ID)
		if err != nil {
			return err
		}

		for _, evDTO := range dto.Events {
			ev, err := s.mToEvent.ToEvent(evDTO)
			if err != nil {
				return err
			}

			if err := v.Apply(ev); err != nil {
				return err
			}
		}

		return s.rVaultPersister.Persist(ctx, v)
	})
}

func (s *service) Pull(ctx context.Context, dto PullDTO) ([]EventDTO, error) {
	if dto.ID == "" {
		return nil, fmt.Errorf("invalid id")
	}

	var events []EventDTO

	err := tx.Run(ctx, func(ctx context.Context) error {
		v, err := s.rVaultGetter.Get(ctx, dto.ID)
		if err != nil {
			return err
		}

		for _, ev := range v.Events(dto.Version) {
			dto, err := s.mToDTO.ToDTO(ev)
			if err != nil {
				return err
			}
			events = append(events, dto)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return events, nil
}
