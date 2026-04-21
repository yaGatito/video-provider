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

	c, err := config.LoadConfig("video")
	if err != nil {
		return fmt.Errorf("failed to load service config: %w", err)
	}

	pgConfig, err := pgxpool.ParseConfig(dbURL(c))
	if err != nil {
		return fmt.Errorf("failed to load service config: %w", err)
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

	mwLog := middleware.NewMiddlewareLogger(os.Stdout, "[VIDSVC]")

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
		auth.NewAuthorizer([]byte(c.JwtSecret)).Auth,
		mwLog.LoggingMiddleware,
		middleware.CORSMiddleware,
	)

	log.Printf("Video-service starting on port %s", c.ApiPort)
	err = http.ListenAndServe(":"+c.ApiPort, router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func dbURL(c config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&pool_max_conns=%d&pool_max_conn_lifetime=%s",
		c.DbUser,
		c.DbPass,
		c.DbHost,
		c.DbPort,
		c.DbName,
		c.ApiSslModCon,
		c.ApiMaxDbCons,
		c.ApiMaxDbConLife)
}
