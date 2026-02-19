package logger

import (
	"io"
	"log"
	"net/http"
	"time"
)

type MiddlewareLogger struct {
	Log *log.Logger
}

func NewMiddlewareLogger(out io.Writer, tag string) *MiddlewareLogger {
	return &MiddlewareLogger{
		Log: log.New(out, tag, log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC),
	}
}

func (l *MiddlewareLogger) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Log.Printf("REQUEST: [%s] %s \"%s\"\n", time.Now().String(), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
