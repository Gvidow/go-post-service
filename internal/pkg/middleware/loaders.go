package middleware

import (
	"context"
	"net/http"

	"github.com/gvidow/go-post-service/internal/pkg/loader"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
)

func WithLoaders(repo usecase.Repository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), loader.LoadersKey, loader.NewLoaders(repo)))
		next.ServeHTTP(w, r)
	})
}
