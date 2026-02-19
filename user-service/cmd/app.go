package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	http2 "user-service/internal/adapters/http"
	"user-service/internal/adapters/mysql"
	usecase "user-service/internal/app"
	logger "video-provider/pkg/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
}

func run() error {
	db, err := sql.Open(
		"mysql", "root:root@tcp(localhost:3306)/user-records")
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %v", err)
	}
	log.Printf("Connected to the database successfully")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to close the connection to the database: %v", err)
		}
	}(db)

	userRepository := mysql.NewSQLUserRepository(db)
	userInteractor := usecase.NewUserService(userRepository)
	userHandler := http2.NewUserHandler(userInteractor)

	mwLog := logger.NewMiddlewareLogger(os.Stdout, "[USRSVC]")
	router := mux.NewRouter()
	router.Use(mwLog.LoggingMiddleware)
	router.HandleFunc("/v1/users", userHandler.Create).Methods("POST")
	router.HandleFunc("/v1/users/{id}", userHandler.Get).Methods("GET")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return fmt.Errorf("Failed to start the server: %v", err)
	}
	fmt.Printf("Server successfully started")
	return nil
}
