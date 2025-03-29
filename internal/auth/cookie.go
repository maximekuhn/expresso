package auth

import (
	"net/http"
	"time"
)

type CookieProvider interface {
	Generate(sessionId string, expiresAt time.Time) http.Cookie
}

const CookieName string = "expresso-delicious-cookie"

type LocalhostCookieProvider struct{}

func NewLocalhostCookieProvider() *LocalhostCookieProvider {
	return &LocalhostCookieProvider{}
}

func (_ *LocalhostCookieProvider) Generate(sessionId string, expiresAt time.Time) http.Cookie {
	return http.Cookie{
		Name:     CookieName,
		Value:    sessionId,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		Secure:   false, // Safari needs it to be false for localhost
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

type ProductionCookieProvider struct {
	domain string
}

func NewProductionCookieProvider(domain string) *ProductionCookieProvider {
	return &ProductionCookieProvider{
		domain: domain,
	}
}

func (pp *ProductionCookieProvider) Generate(sessionId string, expiresAt time.Time) http.Cookie {
	return http.Cookie{
		Name:     CookieName,
		Value:    sessionId,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   pp.domain,
	}
}
