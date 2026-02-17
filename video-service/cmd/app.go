package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	logger "video-provider/pkg/middleware"
	httpadapter "video-service/internal/adapters/http"
	"video-service/internal/adapters/idgen"
	"video-service/internal/adapters/postgres"
	"video-service/internal/app"

	"github.com/joho/godotenv"

	_ "video-service/docs"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
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

// TODO: Fix the problem 1: application should beign configured from one place. Or at least it should be separated.
// TODO: Decide who should be responsible for migration. (the one who run service or the service itself)

func run() error {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}
	connString := os.Getenv("DATABASE_URL")
	port := os.Getenv("API_PORT")

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}
	config.MaxConns = 30
	config.HealthCheckPeriod = time.Minute * 90
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	defer pool.Close()

	videoRepository := postgres.NewVideoRepoPostgreSQL(pool)

	idGen := idgen.New()
	mwLog := logger.NewMiddlewareLogger(os.Stdout, "[VIDSVC]")

	videoService := app.NewVideoInteractor(videoRepository)
	videoHandler := httpadapter.NewVideoHandler(videoService, idGen, mwLog.Log)

	router := mux.NewRouter()
	router.Use(mwLog.LoggingMiddleware)
	httpadapter.SetupRouter(router, videoHandler)

	mwLog.Log.Printf("Server successfully started")
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
