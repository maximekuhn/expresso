package user

import (
	"context"

	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/maximekuhn/expresso/internal/transaction"
	"github.com/maximekuhn/expresso/internal/user"
)

type RegisterUseCaseRequest struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}

type RegisterUseCaseHandler struct {
	sessionProvider transaction.SessionProvider
	authService     *auth.Service
	userService     *user.Service
}

func NewRegisterUseCaseHandler(
	sessionProvider transaction.SessionProvider,
	authService *auth.Service,
	userService *user.Service,
) *RegisterUseCaseHandler {
	return &RegisterUseCaseHandler{
		sessionProvider: sessionProvider,
		authService:     authService,
		userService:     userService,
	}
}

func (h *RegisterUseCaseHandler) Handle(ctx context.Context, r *RegisterUseCaseRequest) error {
	if err := h.validateRequest(r); err != nil {
		return err
	}

	if r.Password != r.PasswordConfirm {
		return PasswordAndConfirmationDontMatchError{}
	}

	session := h.sessionProvider.Provide(ctx)
	return session.Transaction(ctx, func(ctx context.Context) error {
		// create user
		userID, err := h.userService.CreateUser(ctx, r.Name)
		if err != nil {
			return err
		}

		// create auth entry
		return h.authService.CreateAuthEntry(ctx, r.Email, r.Password, userID)
	})
}

func (_ *RegisterUseCaseHandler) validateRequest(r *RegisterUseCaseRequest) error {
	if err := auth.ValidateEmail(r.Email); err != nil {
		return err
	}
	if err := auth.ValidatePassword(r.Password); err != nil {
		return err
	}
	if err := user.ValidateName(r.Name); err != nil {
		return err
	}
	return nil
}
