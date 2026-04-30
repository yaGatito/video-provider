package middleware

import (
	"net/http"
	"video-provider/common/shared"
)

type MiddlewareLogger struct {
	Log *shared.Logger
}

func NewMiddlewareLogger(log *shared.Logger) *MiddlewareLogger {
	return &MiddlewareLogger{
		Log: log,
	}
}

func (l *MiddlewareLogger) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Log.LogRequest(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
