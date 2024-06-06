package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gvidow/go-post-service/pkg/logger"
)

func WithLogger(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLogger := log

		key, ok := r.Context().Value(ReqID).(string)
		if ok {
			reqLogger = reqLogger.WithFields(logger.String("request-id", key))
			r = r.WithContext(context.WithValue(r.Context(), Logger, reqLogger))
		}
		reqLogger.Info("request start", logger.String("path", r.URL.Path))

		defer func(t time.Time) {
			reqLogger.Info("request finish", logger.String("time", strconv.Itoa(int(time.Since(t)))))
		}(time.Now())
		next.ServeHTTP(w, r)
	})
}
