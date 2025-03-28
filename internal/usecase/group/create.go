package group

import (
	"context"

	"github.com/maximekuhn/expresso/internal/group"
	"github.com/maximekuhn/expresso/internal/transaction"
	"github.com/maximekuhn/expresso/internal/user"
)

type CreateUseCaseRequest struct {
	Owner     *user.User
	GroupName string
	Password  string
}

type CreateUseCaseRequestHandler struct {
	sessionProvider transaction.SessionProvider
	groupService    *group.Service
}

func NewCreateUseCaseRequestHandler(
	sessionProvider transaction.SessionProvider,
	groupService *group.Service,
) *CreateUseCaseRequestHandler {
	return &CreateUseCaseRequestHandler{
		sessionProvider: sessionProvider,
		groupService:    groupService,
	}
}

func (h *CreateUseCaseRequestHandler) Handle(ctx context.Context, r *CreateUseCaseRequest) error {
	if err := h.validateRequest(r); err != nil {
		return err
	}

	session := h.sessionProvider.Provide(ctx)
	return session.Transaction(ctx, func(ctx context.Context) error {
		return h.groupService.CreateGroup(
			ctx,
			r.Owner.ID,
			r.GroupName,
			r.Password,
		)
	})
}

func (_ *CreateUseCaseRequestHandler) validateRequest(r *CreateUseCaseRequest) error {
	if err := group.ValidatePassword(r.Password); err != nil {
		return err
	}
	return group.ValidateName(r.GroupName)
}
