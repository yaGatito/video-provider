package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/echo", echoHandler)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      authMiddleware(loggingMiddleware(mux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("starting http server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Graceful shutdown on interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutdown signal received, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("server stopped gracefully")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{"message": "ninja http server", "path": r.URL.Path}
	_ = json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := map[string]interface{}{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(status)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := map[string]string{}
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			query[k] = v[0]
		}
	}
	resp := map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
		"query":  query,
		"proto":  r.Proto,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("x-user-id")
		if userId == "" {
			log.Printf("[%s] %s error user id is not provided\n", r.Method, r.RequestURI)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, CTX_USER_ID, userId)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		// token := r.Header.Get("Authorization")
	})
}

type CTX_KEY string

const (
	CTX_USER_ID CTX_KEY = "user_id"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idFromCtx := r.Context().Value(CTX_USER_ID)
		userId, ok := idFromCtx.(string)
		if !ok {
			log.Printf("[DEBUG LEVEL] User ID not found in context\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		start := time.Now()
		ww := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)
		duration := time.Since(start)
		log.Printf("%s %s [%s] %s %d %s %s\n", start.Format("2006-01-02 15:04:05.99"), userId, r.Method, r.URL.Path, ww.status, duration.String(), formatSizeHeader(r))
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func formatSizeHeader(r *http.Request) string {
	if l := r.Header.Get("Content-Length"); l != "" {
		if _, err := strconv.Atoi(l); err == nil {
			return "len=" + l
		}
	}
	return ""
}
