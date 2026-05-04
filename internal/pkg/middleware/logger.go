package middleware

import (
	"net/http"
	"video-provider/pkg/common"
)

type MiddlewareLogger struct {
	Log *common.Logger
}

func NewMiddlewareLogger(log *common.Logger) *MiddlewareLogger {
	return &MiddlewareLogger{
		Log: log,
	}
}

func (l *MiddlewareLogger) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Log.Info(r.Method + " \"" + r.RequestURI + "\"")
		next.ServeHTTP(w, r)
	})
}
