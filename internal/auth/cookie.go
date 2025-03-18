package auth

import (
	"net/http"
	"time"
)

// cookie := http.Cookie{
// 	Name:     CookieName,
// 	Value:    sessionId,
// 	MaxAge:   int(time.Until(cookieExpiryDate).Seconds()),
// 	Secure:   true,
// 	HttpOnly: true,
// 	SameSite: http.SameSiteStrictMode,
// }
//

const CookieName string = "expresso-delicious-cookie"

func GenerateCookie(sessionId string, expiresAt time.Time) http.Cookie {
	// TODO: add domain
	return http.Cookie{
		Name:     CookieName,
		Value:    sessionId,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}
