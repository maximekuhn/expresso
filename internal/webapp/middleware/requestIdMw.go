package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type RequestIdKey struct{}

type RequestIdMiddleware struct{}

func NewRequestIdMiddleware() *RequestIdMiddleware {
	return &RequestIdMiddleware{}
}

func (mw *RequestIdMiddleware) Middleware(next http.Handler) http.Handler {
	const header = "X-Request-Id"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId, ok := r.Context().Value(RequestIdKey{}).(string)
		if !ok || reqId == "" {
			reqId = uuid.NewString()
		}
		ctx := context.WithValue(r.Context(), RequestIdKey{}, reqId)
		w.Header().Set(header, reqId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
