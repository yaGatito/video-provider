package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "video-provider/docs"
	config "video-provider/internal/pkg/config"
	"video-provider/internal/pkg/middleware"
	httpadp "video-provider/internal/video-service/adapters/http"
	"video-provider/internal/video-service/adapters/postgres"
	"video-provider/internal/video-service/app"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const configName = "video-service"

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
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg, err := config.ParseFromFS(configName)
	if err != nil {
		return fmt.Errorf("failed to parse config bytes: %w", err)
	}

	pgConfig, err := pgxpool.ParseConfig(cfg.Db.GetURL())
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}
	// cfg.MaxConns = 30
	// cfg.HealthCheckPeriod = time.Minute * 90
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	defer pool.Close()

	mwLog := middleware.NewMiddlewareLogger(os.Stdout, "[VIDSVC]")

	videoRepository := postgres.NewVideoRepoPostgreSQL(pool)
	videoService := app.NewVideoInteractor(videoRepository)
	videoHandler := httpadp.NewVideoHandler(videoService, mwLog.Log)

	router := mux.NewRouter()
	router.Use(middleware.CORSMiddleware)
	router.Use(mwLog.LoggingMiddleware)

	httpadp.SetupRouter(router, videoHandler)

	log.Printf("Video-service starting on port %s", cfg.Api.Port)
	err = http.ListenAndServe(":"+cfg.Api.Port, router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
