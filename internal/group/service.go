package group

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

func (s *Service) CreateGroup(
	ctx context.Context,
	ownerID uuid.UUID,
	name string,
	password string,
) error {
	if err := ValidatePassword(password); err != nil {
		return err
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	groupID := s.idProvider.Provide()
	createdAt := s.datetimeProvider.Provide()
	g, err := New(
		groupID,
		name,
		ownerID,
		make([]uuid.UUID, 0),
		hashedPassword,
		createdAt,
	)
	if err != nil {
		return err
	}
	return s.store.Save(ctx, *g)
}

func (s *Service) ListGroupOfUser(ctx context.Context, userID uuid.UUID) ([]Group, error) {
	groupsAsOwner, err := s.store.GetAllWhereUserIsOwner(ctx, userID)
	if err != nil {
		return nil, err
	}
	groupsAsMember, err := s.store.GetAllWhereUserIsMember(ctx, userID)
	if err != nil {
		return nil, err
	}
	return append(groupsAsOwner, groupsAsMember...), nil
}
