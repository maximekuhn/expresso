package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/maximekuhn/expresso/internal/logger"
	"github.com/maximekuhn/expresso/internal/transaction"
	"github.com/maximekuhn/expresso/internal/user"
)

type UserKey struct{}

// LoggedInMiddleware checks if an auth cookie is found and if the associated
// session is valid. If none of these conditions is true, the caller is redirected
// to the login page.
// The found user is injected in the request context.
type LoggedInMiddleware struct {
	logger          *slog.Logger
	authService     *auth.Service
	userService     *user.Service
	sessionProvider transaction.SessionProvider
}

func NewLoggedInMiddleware(
	l *slog.Logger,
	authService *auth.Service,
	userService *user.Service,
	sessionProvider transaction.SessionProvider,
) *LoggedInMiddleware {
	return &LoggedInMiddleware{
		logger:          l.With(slog.String(logger.LoggerNameField, "LoggedInMiddleware")),
		authService:     authService,
		userService:     userService,
		sessionProvider: sessionProvider,
	}
}

func (mw *LoggedInMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logger.UpgradeWithRequestId(r.Context(), RequestIdKey{}, mw.logger)
		cookie, err := r.Cookie(auth.CookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				l.Info("no cookie")
			} else {
				l.Error("failed to retrieve cookie", slog.String("err", err.Error()))
			}

			// redirect to login page and don't call next handler
			w.Header().Add("HX-Redirect", "/login") // for HTMX callers
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var u *user.User

		ctx := r.Context()
		if err := mw.sessionProvider.Provide(ctx).Transaction(ctx, func(ctx context.Context) error {
			userID, valid, err := mw.authService.IsSessionValid(ctx, cookie)
			if err != nil {
				l.Error("failed to check session validity", slog.String("err", err.Error()))
				return err
			}
			if !valid {
				l.Info("session is not valid")
				return nil
			}
			usr, found, err := mw.userService.Get(ctx, userID)
			if err != nil {
				l.Error("failed to retrieve user", slog.String("err", err.Error()))
				return err
			}
			if !found {
				l.Info("user not found", slog.String("userId", userID.String()))
				return nil
			}
			u = usr
			return nil
		}); err != nil {
			l.Error("failed to run transaction", slog.String("err", err.Error()))
			return
		}

		if u == nil {
			// no user found or no valid session
			return
		}

		// inject user in the request context to be accessed
		// by next middleware(s)/handler(s)
		updatedCtx := context.WithValue(ctx, UserKey{}, u)
		next.ServeHTTP(w, r.WithContext(updatedCtx))
	})

}
