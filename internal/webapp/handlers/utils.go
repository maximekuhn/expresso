package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/expresso/internal/user"
	"github.com/maximekuhn/expresso/internal/webapp/middleware"
	uierrors "github.com/maximekuhn/expresso/internal/webapp/ui/components/errors"
)

func returnBadRequestAndBoxError(errMsg string, l *slog.Logger, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	boxError(errMsg, l, w, r)
}

func returnConflictAndBoxError(errMsg string, l *slog.Logger, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusConflict)
	boxError(errMsg, l, w, r)
}

func boxError(errMsg string, l *slog.Logger, w http.ResponseWriter, r *http.Request) {
	if err := uierrors.BoxError(errMsg).Render(r.Context(), w); err != nil {
		l.Error(
			"failed to render errors.BoxError",
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func extractUserOrReturnInternalError(l *slog.Logger, w http.ResponseWriter, r *http.Request) *user.User {
	u, ok := r.Context().Value(middleware.UserKey{}).(*user.User)
	if !ok {
		l.Error("found data in context with middleware.UserKey but could not cast")
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	return u
}
