package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/expresso/internal/logger"
	"github.com/maximekuhn/expresso/internal/webapp/middleware"
	"github.com/maximekuhn/expresso/internal/webapp/ui/templates/pages"
)

type IndexHandler struct {
	logger *slog.Logger
}

func NewIndexHandler(l *slog.Logger) *IndexHandler {
	return &IndexHandler{
		logger: l.With(slog.String(logger.LoggerNameField, "IndexHandler")),
	}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getIndexPage(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *IndexHandler) getIndexPage(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeWithRequestId(r.Context(), middleware.RequestIdKey{}, h.logger)
	if err := pages.Index().Render(r.Context(), w); err != nil {
		l.Error(
			"failed to render pages.Index",
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
