package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ReqID, reqID)))
	})
}
