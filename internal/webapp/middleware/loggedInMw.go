package middleware

import (
	"context"
	"net/http"

	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/maximekuhn/expresso/internal/transaction"
	"github.com/maximekuhn/expresso/internal/user"
)

type UserKey struct{}

// LoggedInMiddleware checks if an auth cookie is found and if the associated
// session is valid. If none of these conditions is true, the caller is redirected
// to the login page.
// The found user is injected in the request context.
type LoggedInMiddleware struct {
	authService     *auth.Service
	userService     *user.Service
	sessionProvider transaction.SessionProvider
}

func NewLoggedInMiddleware(
	authService *auth.Service,
	userService *user.Service,
	sessionProvider transaction.SessionProvider,
) *LoggedInMiddleware {
	return &LoggedInMiddleware{
		authService:     authService,
		userService:     userService,
		sessionProvider: sessionProvider,
	}
}

func (mw *LoggedInMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(auth.CookieName)
		if err != nil {
			// redirect to login page and don't call next handler
			w.Header().Add("HX-Redirect", "/login") // for HTMX callers
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var u *user.User

		ctx := r.Context()
		mw.sessionProvider.Provide(ctx).Transaction(ctx, func(ctx context.Context) error {
			userID, valid, err := mw.authService.IsSessionValid(ctx, cookie)
			if err != nil {
				return err
			}
			if !valid {
				return nil
			}
			usr, found, err := mw.userService.Get(ctx, userID)
			if err != nil {
				return err
			}
			if !found {
				return nil
			}
			u = usr
			return nil
		})

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
