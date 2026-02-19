package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	http2 "video-provider/internal/user-service/adapters/http"
	"video-provider/internal/user-service/adapters/postgres"
	usecase "video-provider/internal/user-service/app"
	logger "video-provider/pkg/middleware"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		// Default to localhost:5433 (mapped in Makefile/Docker)
		dbUrl = "postgres://gato:root@localhost:5433/userdb?sslmode=disable"
	}

	dbPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		return fmt.Errorf("Failed to ping the database: %v", err)
	}
	log.Printf("Connected to the database successfully")

	userRepository := postgres.NewPostgresUserRepository(dbPool)
	userInteractor := usecase.NewUserService(userRepository)
	userHandler := http2.NewUserHandler(userInteractor)

	mwLog := logger.NewMiddlewareLogger(os.Stdout, "[USRSVC]")
	router := mux.NewRouter()
	router.Use(mwLog.LoggingMiddleware)
	router.HandleFunc("/v1/users", userHandler.Create).Methods("POST")
	router.HandleFunc("/v1/users/{id}", userHandler.Get).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		return fmt.Errorf("Failed to start the server: %v", err)
	}
	fmt.Printf("Server successfully started")
	return nil
}
