package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "user-service/docs"

	"user-service/pkg/auth"
	"user-service/pkg/middleware"

	cryptoadp "user-service/adapters/crypto"
	httpadp "user-service/adapters/http"
	"user-service/adapters/postgres"
	"user-service/app"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbUser  = "POSTGRES_USER"
	dbPass  = "POSTGRES_PASSWORD"
	dbHost  = "USER_DB_HOST"
	dbPort  = "USER_DB_PORT"
	dbName  = "USER_DB_NAME"
	apiPort = "USER_API_PORT"
)

// @title			User Service API
// @version			1.0
// @description		Service for managing users.
// @host			localhost:8081
// @BasePath		/
func main() {
	if err := run(); err != nil {
		log.Println(err)
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
	dbPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	err = dbPool.Ping(ctx)
	if err != nil {
		log.Default().Printf("failed to ping database: %s\n", err.Error())
	}

	defer dbPool.Close()

	mwLog := middleware.NewMiddlewareLogger(os.Stdout, "[USRSVC]")

	userRepository := postgres.NewPostgresUserRepository(dbPool)
	pwHasher := cryptoadp.NewBCryptPasswordHasher()
	userInteractor := app.NewUserService(userRepository, pwHasher, auth.GetJWTSecret)
	userHandler := httpadp.NewUserHandler(userInteractor, mwLog.Log)

	router := mux.NewRouter()

	httpadp.SetupRouter(
		router,
		userHandler,
		auth.Auth,
		middleware.CORSMiddleware,
		mwLog.LoggingMiddleware,
	)

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
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=30&pool_max_conn_lifetime=1h30m",
		os.Getenv(
			dbUser,
		),
		os.Getenv(dbPass),
		os.Getenv(dbHost),
		os.Getenv(dbPort),
		os.Getenv(dbName),
	)
}
