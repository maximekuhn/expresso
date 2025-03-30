package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/maximekuhn/expresso/internal/logger"
	usecaseUser "github.com/maximekuhn/expresso/internal/usecase/user"
	"github.com/maximekuhn/expresso/internal/webapp/middleware"
)

type LogoutHandler struct {
	logger         *slog.Logger
	logoutUseCase  *usecaseUser.LogoutUseCaseHandler
	cookieProvider auth.CookieProvider
}

func NewLogoutHandler(
	l *slog.Logger,
	logoutUseCase *usecaseUser.LogoutUseCaseHandler,
	cookieProvider auth.CookieProvider,
) *LogoutHandler {
	return &LogoutHandler{
		logger:         l.With(slog.String(logger.LoggerNameField, "LogoutHandler")),
		logoutUseCase:  logoutUseCase,
		cookieProvider: cookieProvider,
	}
}
func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.logout(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *LogoutHandler) logout(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeWithRequestId(r.Context(), middleware.RequestIdKey{}, h.logger)

	loggedUser := extractUserOrReturnInternalError(l, w, r)
	if loggedUser == nil {
		return
	}

	if err := h.logoutUseCase.Handle(r.Context(), &usecaseUser.LogoutUseCaseRequest{
		User: loggedUser,
	}); err != nil {
		h.handleLogoutError(l, w, r, err)
		return
	}

	w.Header().Add("HX-Redirect", "/login")
	cookie := h.cookieProvider.Generate("expired", time.Now().UTC())
	http.SetCookie(w, &cookie)
	l.Info("user logged out successfully")
}

func (_ *LogoutHandler) handleLogoutError(l *slog.Logger, w http.ResponseWriter, _ *http.Request, err error) {
	var notFoundError auth.EntryNotFoundError
	if errors.As(err, &notFoundError) {
		l.Info("no session found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	l.Error("internal error", slog.String("err", err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}
