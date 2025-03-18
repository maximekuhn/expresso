package handlers

import (
	"log/slog"
	"net/http"

	uierrors "github.com/maximekuhn/expresso/internal/webapp/ui/components/errors"
)

func returnBadRequestAndBoxError(errMsg string, l *slog.Logger, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	if err := uierrors.BoxError(errMsg).Render(r.Context(), w); err != nil {
		l.Error(
			"failed to render errors.BoxError",
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
