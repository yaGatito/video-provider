package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	httpadp "video-service/adapters/http"
	"video-service/adapters/postgres"
	"video-service/app"
	_ "video-service/docs"
	"video-service/pkg/auth"
	"video-service/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbUser       = "POSTGRES_USER"
	dbPass       = "POSTGRES_PASSWORD"
	dbHost       = "VIDEO_DB_HOST"
	dbPort       = "VIDEO_DB_PORT"
	dbName       = "VIDEO_DB_NAME"
	apiPort      = "VIDEO_API_PORT"
)

// @title			Video Service API
// @version		1.0
// @description	Service for managing video content.
// @host			localhost:8080
// @BasePath		/
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	pgConfig, err := pgxpool.ParseConfig(dbUrl())
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}
	// cfg.MaxConns = 30
	// cfg.HealthCheckPeriod = time.Minute * 90
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Default().Printf("failed to ping database: %s\n", err.Error())
	}
	defer pool.Close()

	mwLog := middleware.NewMiddlewareLogger(os.Stdout, "[VIDSVC]")

	videoRepository := postgres.NewVideoRepoPostgreSQL(pool)
	videoService := app.NewVideoInteractor(videoRepository)
	videoHandler := httpadp.NewVideoHandler(videoService, mwLog.Log)

	router := mux.NewRouter()

	httpadp.SetupRouter(router, &videoHandler, auth.Auth, mwLog.LoggingMiddleware, middleware.CORSMiddleware)

	log.Printf("Video-service starting on port %s", os.Getenv(apiPort))
	err = http.ListenAndServe(":"+os.Getenv(apiPort), router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// dbUrl must be called only after setup OS env variables.
func dbUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=30&pool_max_conn_lifetime=1h30m",
		os.Getenv(dbUser), os.Getenv(dbPass), os.Getenv(dbHost), os.Getenv(dbPort), os.Getenv(dbName))
}
