package group

import (
	"context"

	"github.com/maximekuhn/expresso/internal/group"
	"github.com/maximekuhn/expresso/internal/transaction"
	"github.com/maximekuhn/expresso/internal/user"
)

type JoinUseCaseRequest struct {
	User      *user.User
	GroupName string
	Password  string
}

type JoinUseCaseRequestHandler struct {
	sessionProvider transaction.SessionProvider
	groupService    *group.Service
}

func NewJoinUseCaseRequestHandler(
	sessionProvider transaction.SessionProvider,
	groupService *group.Service,
) *JoinUseCaseRequestHandler {
	return &JoinUseCaseRequestHandler{
		sessionProvider: sessionProvider,
		groupService:    groupService,
	}
}

func (h *JoinUseCaseRequestHandler) Handle(ctx context.Context, r *JoinUseCaseRequest) error {
	if err := h.validateRequest(r); err != nil {
		return err
	}
	session := h.sessionProvider.Provide(ctx)
	return session.Transaction(ctx, func(ctx context.Context) error {
		return h.groupService.JoinGroup(ctx, r.User.ID, r.GroupName, r.Password)
	})
}

func (_ *JoinUseCaseRequestHandler) validateRequest(r *JoinUseCaseRequest) error {
	if err := group.ValidateName(r.GroupName); err != nil {
		return err
	}
	return group.ValidatePassword(r.Password)
}
