package handlers

import (
	"log/slog"
	"net/http"

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
		// TODO: check error type and handle it properly
		l.Error("failed to register user", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
