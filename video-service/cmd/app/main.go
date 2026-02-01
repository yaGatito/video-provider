package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	httpadapter "video-service/internal/adapters/http"
	"video-service/internal/adapters/idgen"
	"video-service/internal/adapters/postgres"
	"video-service/internal/app"

	"github.com/joho/godotenv"

	_ "video-service/docs"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

// @title           Video Service API
// @version         1.0
// @description     Service for managing video content.
// @host            localhost:8080
// @BasePath        /
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connString := os.Getenv("DATABASE_URL")
	port := os.Getenv("API_PORT")

	config, err := pgxpool.ParseConfig(connString)
	config.MaxConns = 30
	if err != nil {
		log.Fatal(err)
	}
	config.HealthCheckPeriod = time.Minute * 90
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	videoRepository := postgres.NewVideoRepoPostgreSQL(pool)

	idGen := idgen.New()
	mwLog := MiddlewareLogger{
		log: log.New(os.Stdout, "[VIDSVC] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC),
	}

	videoService := app.NewVideoInteractor(videoRepository)
	videoHandler := httpadapter.NewVideoHandler(videoService, idGen, mwLog.log)

	router := mux.NewRouter()
	router.Use(mwLog.loggingMiddleware)
	httpadapter.SetupRouter(router, videoHandler)

	mwLog.log.Printf("Server successfully started")
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		return err
	}
	return nil
}

type MiddlewareLogger struct {
	log *log.Logger
}

func (l *MiddlewareLogger) Log() *log.Logger {
	return l.log
}

func (l *MiddlewareLogger) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: [%s] %s \"%s\"\n", time.Now().String(), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
