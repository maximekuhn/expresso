package user

import (
	"context"

	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/maximekuhn/expresso/internal/transaction"
	"github.com/maximekuhn/expresso/internal/user"
)

type LogoutUseCaseRequest struct {
	User *user.User
}

type LogoutUseCaseHandler struct {
	sessionProvider transaction.SessionProvider
	authService     *auth.Service
}

func NewLogoutUseCaseHandler(
	sessionProvider transaction.SessionProvider,
	authService *auth.Service,
) *LogoutUseCaseHandler {
	return &LogoutUseCaseHandler{
		sessionProvider: sessionProvider,
		authService:     authService,
	}
}

func (h *LogoutUseCaseHandler) Handle(ctx context.Context, r *LogoutUseCaseRequest) error {
	return h.sessionProvider.Provide(ctx).Transaction(ctx, func(ctx context.Context) error {
		return h.authService.RevokeSession(ctx, r.User.ID)
	})
}
