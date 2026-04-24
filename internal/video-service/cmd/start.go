package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"video-provider/common/auth"
	"video-provider/common/config"
	httpadp "video-provider/video-service/adapters/http"
	"video-provider/video-service/adapters/postgres"
	"video-provider/video-service/app"
	_ "video-provider/video-service/docs"

	"video-provider/common/middleware"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

// @title			Video Service API
// @version			1.0
// @description		Service for managing video content.
// @host			localhost:8080
// @BasePath		/
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	jsonConfig, err := os.ReadFile("./config/video_config.json")
	if err != nil {
		return fmt.Errorf("failed to load json service config: %w", err)
	}
	c, err := config.LoadConfig("video", jsonConfig)
	if err != nil {
		return fmt.Errorf("failed to load service config: %w", err)
	}

	pgConfig, err := pgxpool.ParseConfig(dbURL(c))
	if err != nil {
		return fmt.Errorf("failed to parse db config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Default().Printf("failed to ping database: %s\n", err.Error())
	}
	defer pool.Close()

	mwLog := middleware.NewMiddlewareLogger(httpadp.DefaultLogger)

	authorizer := auth.NewAuthorizer(auth.NewTokenizer(c))
	videoRepository := postgres.NewVideoRepoPostgreSQL(pool)
	videoService := app.NewVideoInteractor(videoRepository)

	val, err := httpadp.NewVideoValidator()
	if err != nil {
		log.Default().Printf("failed to create validator: %s\n", err.Error())
	}
	videoHandler := httpadp.NewVideoHandler(videoService, mwLog.Log, val)

	router := mux.NewRouter()

	httpadp.SetupRouter(
		router,
		videoHandler,
		authorizer.Auth,
		mwLog.LoggingMiddleware,
		middleware.CORSMiddleware,
	)

	log.Printf("Video-service starting on port %s", c.EnvConf.ApiPort)
	err = http.ListenAndServe(":"+c.EnvConf.ApiPort, router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
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
