package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"video-provider/pkg/common"
	_ "video-provider/user-service/docs"

	"video-provider/pkg/auth"
	"video-provider/pkg/config"
	"video-provider/pkg/middleware"

	cryptoadp "video-provider/user-service/adapters/crypto"
	httpadp "video-provider/user-service/adapters/http"
	"video-provider/user-service/adapters/postgres"
	"video-provider/user-service/adapters/postgres/sqlcgen"
	"video-provider/user-service/app"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

// @title			User Service API
// @version			1.0
// @description		Service for managing users.
// @host			localhost:8081
// @BasePath		/
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	jsonConfig, err := os.ReadFile("./config/user_config.json")
	if err != nil {
		return fmt.Errorf("failed to load json service config: %w", err)
	}
	c, err := config.LoadConfig("user", jsonConfig)
	if err != nil {
		return fmt.Errorf("failed to load service config: %w", err)
	}

	pgConfig, err := pgxpool.ParseConfig(dbURL(c))
	if err != nil {
		return fmt.Errorf("failed to parse connection config: %w", err)
	}
	dbPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	err = dbPool.Ping(ctx)
	if err != nil {
		log.Default().Printf("failed to ping database: %s\n", err.Error())
	}

	defer dbPool.Close()

	log := common.NewLogger(common.DefaultOutput, "USR-SVC")

	userRepository := postgres.NewPostgresUserRepository(sqlcgen.New(dbPool))
	pwHasher := cryptoadp.NewBCryptPasswordHasher()
	tokenizer := auth.NewTokenizer(c)
	authorizer := auth.NewAuthorizer(tokenizer)
	userInteractor := app.NewUserService(userRepository, pwHasher, tokenizer)
	userHandler := httpadp.NewUserHandler(userInteractor, log)

	router := mux.NewRouter()

	httpadp.SetupRouter(
		router,
		userHandler,
		authorizer.Auth,
		middleware.NewMiddlewareLogger(log).LoggingMiddleware,
		middleware.CORSMiddleware,
	)

	log.Info("User-service starting on port: " + c.EnvConf.ApiPort)
	err = http.ListenAndServe(":"+c.EnvConf.ApiPort, router)
	if err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}
	log.Info("Server successfully started")
	return nil
}

func dbURL(c config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&pool_max_conns=%d&pool_max_conn_lifetime=%s",
		c.EnvConf.DbUser,
		c.EnvConf.DbPass,
		c.EnvConf.DbHost,
		c.EnvConf.DbPort,
		c.EnvConf.DbName,
		c.JsonConf.SSLMode,
		c.JsonConf.PoolCons,
		c.JsonConf.PoolConLifetime)
}
