package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/expresso/internal/common"
)

type Service struct {
	store            Store
	idProvider       common.IdProvider
	datetimeProvider common.DatetimeProvider
}

func NewService(
	store Store,
	idProvider common.IdProvider,
	datetimeProvider common.DatetimeProvider,
) *Service {
	return &Service{
		store:            store,
		idProvider:       idProvider,
		datetimeProvider: datetimeProvider,
	}
}

func (s *Service) CreateUser(ctx context.Context, name string) (uuid.UUID, error) {
	id := s.idProvider.Provide()
	createdAt := s.datetimeProvider.Provide()
	u, err := New(id, name, createdAt)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, s.store.Save(ctx, *u)
}

func (s *Service) Get(ctx context.Context, userID uuid.UUID) (*User, bool, error) {
	return s.store.GetById(ctx, userID)
}
