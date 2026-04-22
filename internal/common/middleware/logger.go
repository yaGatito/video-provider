package middleware

import (
	"log"
	"net/http"
	"time"
)

type MiddlewareLogger struct {
	Log *log.Logger
}

func NewMiddlewareLogger(log *log.Logger) *MiddlewareLogger {
	return &MiddlewareLogger{
		Log: log,
	}
}

func (l *MiddlewareLogger) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Log.Printf("REQUEST: [%s] %s \"%s\"\n", time.Now().String(), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
