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

type LoginHandler struct {
	logger       *slog.Logger
	loginUseCase *usecaseUser.LoginUseCaseHandler
}

func NewLoginHandler(
	l *slog.Logger,
	loginUseCase *usecaseUser.LoginUseCaseHandler,
) *LoginHandler {
	return &LoginHandler{
		logger:       l.With(slog.String(logger.LoggerNameField, "LoginHandler")),
		loginUseCase: loginUseCase,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getLoginPage(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.login(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *LoginHandler) getLoginPage(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeWithRequestId(r.Context(), middleware.RequestIdKey{}, h.logger)
	if err := pages.Login().Render(r.Context(), w); err != nil {
		l.Error(
			"failed to render pages.Login",
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *LoginHandler) login(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeWithRequestId(r.Context(), middleware.RequestIdKey{}, h.logger)

	if err := r.ParseForm(); err != nil {
		l.Error("failed to parse form", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	res, err := h.loginUseCase.Handle(r.Context(), &usecaseUser.LoginUseCaseRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		h.handleError(l, w, r, err)
		return
	}

	cookie := auth.GenerateCookie(res.SessionId, res.ExpiresAt)
	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/")
	l.Info("user logged in successfully")
}

func (_ *LoginHandler) handleError(l *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	var badCredentialsError auth.BadCredentialsError
	if errors.As(err, &badCredentialsError) {
		l.Info("bad credentials", slog.String("err", badCredentialsError.Error()))
		returnBadRequestAndBoxError("Invalid credentials", l, w, r)
		return
	}

	var entryNotFoundError auth.EntryNotFoundError
	if errors.As(err, &entryNotFoundError) {
		l.Info("user not found", slog.String("err", entryNotFoundError.Error()))
		returnBadRequestAndBoxError("Invalid credentials", l, w, r)
		return
	}

	l.Error("internal error", slog.String("err", err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}
