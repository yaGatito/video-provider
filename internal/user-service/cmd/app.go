package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "user-service/docs"

	"github.com/yaGatito/video-provider/internal/pkg/middleware	"

	cryptoadp "user-service/adapters/crypto"
	httpadp "user-service/adapters/http"
	"user-service/adapters/postgres"
	"user-service/app"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	dbUser  = "POSTGRES_USER"
	dbPass  = "POSTGRES_PASSWORD"
	dbHost  = "USER_DB_HOST"
	dbPort  = "USER_DB_PORT"
	dbName  = "USER_DB_NAME"
	apiPort = "API_PORT"
)

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

	pgConfig, err := pgxpool.ParseConfig(dbUrl())
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

	mwLog := middleware.NewMiddlewareLogger(os.Stdout, "[USRSVC]")

	userRepository := postgres.NewPostgresUserRepository(dbPool)
	pwHasher := cryptoadp.NewBCryptPasswordHasher()
	userInteractor := app.NewUserService(userRepository, pwHasher)
	userHandler := httpadp.NewUserHandler(userInteractor, mwLog.Log)

	router := mux.NewRouter()
	router.Use(middleware.CORSMiddleware)
	router.Use(mwLog.LoggingMiddleware)

	httpadp.SetupRouter(router, userHandler)

	log.Printf("User-service starting on port %s", os.Getenv(apiPort))
	err = http.ListenAndServe(":"+os.Getenv(apiPort), router)
	if err != nil {
		return fmt.Errorf("Failed to start the server: %v", err)
	}
	fmt.Printf("Server successfully started")
	return nil
}

// dbUrl must be called only after setup OS env variables.
func dbUrl() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=%s&pool_max_conn_lifetime=1h30m",
		"postgres", os.Getenv(dbUser), os.Getenv(dbPass), os.Getenv(dbHost), os.Getenv(dbPort), os.Getenv(dbName), "30")
}
