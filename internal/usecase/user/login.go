package user

import (
	"context"
	"time"

	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/maximekuhn/expresso/internal/common"
	"github.com/maximekuhn/expresso/internal/transaction"
)

type LoginUseCaseRequest struct {
	Email    string
	Password string
}

type LoginUseCaseResponse struct {
	SessionId string
	ExpiresAt time.Time
}

type LoginUseCaseHandler struct {
	sessionProvider  transaction.SessionProvider
	authService      *auth.Service
	datetimeProvider common.DatetimeProvider
}

func NewLoginUseCaseHandler(
	sessionProvider transaction.SessionProvider,
	authService *auth.Service,
	datetimeProvider common.DatetimeProvider,
) *LoginUseCaseHandler {
	return &LoginUseCaseHandler{
		sessionProvider:  sessionProvider,
		authService:      authService,
		datetimeProvider: datetimeProvider,
	}
}

func (h *LoginUseCaseHandler) Handle(
	ctx context.Context,
	r *LoginUseCaseRequest,
) (*LoginUseCaseResponse, error) {
	if err := h.validateRequest(r); err != nil {
		return nil, err
	}

	var res *LoginUseCaseResponse

	session := h.sessionProvider.Provide(ctx)
	err := session.Transaction(ctx, func(ctx context.Context) error {
		authEntry, err := h.authService.CreateSession(ctx, r.Email, r.Password)
		if err != nil {
			return err
		}
		res = &LoginUseCaseResponse{
			SessionId: authEntry.SessionID,
			ExpiresAt: *authEntry.SessionExpiresAt,
		}
		return nil

	})
	return res, err
}

func (_ *LoginUseCaseHandler) validateRequest(r *LoginUseCaseRequest) error {
	if err := auth.ValidateEmail(r.Email); err != nil {
		return err
	}
	return auth.ValidatePassword(r.Password)
}
