package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	httpadapter "video-provider/internal/video-service/adapters/http"
	logger "video-provider/pkg/middleware"
	"video-service/internal/adapters/idgen"
	"video-service/internal/adapters/postgres"
	"video-service/internal/app"

	"github.com/joho/godotenv"

	_ "video-service/docs"

	config "video-provider/internal/pkg/config"

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

	fileCfg, err := os.ReadFile(config.СonfigPath)
	if err != nil {
		return fmt.Errorf("failed to read file from %s: %w", config.СonfigPath, err)
	}
	cfg, err := config.ParseConfig(fileCfg)
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

	videoRepository := postgres.NewVideoRepoPostgreSQL(pool)

	idGen := idgen.New()
	mwLog := logger.NewMiddlewareLogger(os.Stdout, "[VIDSVC]")

	videoService := app.NewVideoInteractor(videoRepository)
	videoHandler := httpadapter.NewVideoHandler(videoService, idGen, mwLog.Log)

	router := mux.NewRouter()
	router.Use(mwLog.LoggingMiddleware)
	httpadapter.SetupRouter(router, videoHandler)

	mwLog.Log.Printf("Server successfully started")
	err = http.ListenAndServe(":"+cfg.Api.Port, router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
