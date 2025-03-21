package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/maximekuhn/expresso/internal/logger"
	usecaseUser "github.com/maximekuhn/expresso/internal/usecase/user"
	"github.com/maximekuhn/expresso/internal/webapp/middleware"
	"github.com/maximekuhn/expresso/internal/webapp/ui/templates/pages"
)

type RegisterHandler struct {
	logger          *slog.Logger
	registerUseCase *usecaseUser.RegisterUseCaseHandler
}

func NewRegisterHandler(
	l *slog.Logger,
	registerUseCase *usecaseUser.RegisterUseCaseHandler,
) *RegisterHandler {
	return &RegisterHandler{
		logger:          l.With(slog.String(logger.LoggerNameField, "RegisterHandler")),
		registerUseCase: registerUseCase,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getRegisterPage(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.register(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *RegisterHandler) getRegisterPage(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeWithRequestId(r.Context(), middleware.RequestIdKey{}, h.logger)
	if err := pages.Register().Render(r.Context(), w); err != nil {
		l.Error(
			"failed to render pages.Register",
			slog.String("err", err.Error()),
		)
	}
}

func (h *RegisterHandler) register(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeWithRequestId(r.Context(), middleware.RequestIdKey{}, h.logger)
	if err := r.ParseForm(); err != nil {
		l.Error("failed to parse form", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	email := r.PostForm.Get("email")
	name := r.PostForm.Get("name")
	password := r.PostForm.Get("password")
	passwordConfirm := r.PostForm.Get("password-confirm")

	if err := h.registerUseCase.Handle(r.Context(), &usecaseUser.RegisterUseCaseRequest{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}); err != nil {
		h.handleError(l, w, r, err)
		return
	}

	w.Header().Add("HX-Redirect", "/login")
	l.Info("user registered successfully")
}

func (_ *RegisterHandler) handleError(l *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	var passwordAndConfirmationDontMatchError usecaseUser.PasswordAndConfirmationDontMatchError
	if errors.As(err, &passwordAndConfirmationDontMatchError) {
		l.Info("password and confirmation don't match")
		returnBadRequestAndBoxError("Password and confirmation must match", l, w, r)
		return
	}

	var entryAlreadyExistsError auth.EntryAlreadyExistsError
	if errors.As(err, &entryAlreadyExistsError) {
		l.Info("email already taken")
		returnBadRequestAndBoxError("This email is not available. Try another one", l, w, r)
		return
	}

	l.Error("internal error", slog.String("err", err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}
