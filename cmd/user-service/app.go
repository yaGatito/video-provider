package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "video-provider/docs"
	"video-provider/internal/pkg/config"
	logger "video-provider/internal/pkg/middleware"

	cryptoadp "video-provider/internal/user-service/adapters/crypto"
	httpadp "video-provider/internal/user-service/adapters/http"
	"video-provider/internal/user-service/adapters/postgres"
	"video-provider/internal/user-service/app"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const configName = "user-service"

// @title           User Service API
// @version         1.0
// @description     Service for managing users.
// @host            localhost:8081
// @BasePath        /
func main() {
	if err := run(); err != nil {
		log.Println(err)
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
	dbPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	defer dbPool.Close()

	mwLog := logger.NewMiddlewareLogger(os.Stdout, "[USRSVC]")

	userRepository := postgres.NewPostgresUserRepository(dbPool)
	pwHasher := cryptoadp.NewBCryptPasswordHasher()
	userInteractor := app.NewUserService(userRepository, pwHasher)
	userHandler := httpadp.NewUserHandler(userInteractor, mwLog.Log)

	router := mux.NewRouter()
	router.Use(logger.CORSMiddleware)
	router.Use(mwLog.LoggingMiddleware)

	httpadp.SetupRouter(router, userHandler)

	log.Printf("User-service starting on port %s", cfg.Api.Port)
	err = http.ListenAndServe(":"+cfg.Api.Port, router)
	if err != nil {
		return fmt.Errorf("Failed to start the server: %v", err)
	}
	fmt.Printf("Server successfully started")
	return nil
}
