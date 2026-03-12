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
	httpadp "video-provider/internal/user-service/adapters/http"
	"video-provider/internal/user-service/adapters/postgres"
	usecase "video-provider/internal/user-service/app"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
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
	userInteractor := usecase.NewUserService(userRepository)
	userHandler := httpadp.NewUserHandler(userInteractor)

	router := mux.NewRouter()
	router.Use(logger.CORSMiddleware)
	router.Use(mwLog.LoggingMiddleware)

	router.HandleFunc("/v1/users", userHandler.CreateUser).
		Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/v1/users/{id}", userHandler.GetUser).
		Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/login", userHandler.Login).
		Methods(http.MethodPost, http.MethodOptions)

	router.PathPrefix("/v1/swagger/").HandlerFunc(httpSwagger.WrapHandler)

	log.Printf("User-service starting on port %s", cfg.Api.Port)
	err = http.ListenAndServe(":"+cfg.Api.Port, router)
	if err != nil {
		return fmt.Errorf("Failed to start the server: %v", err)
	}
	fmt.Printf("Server successfully started")
	return nil
}
