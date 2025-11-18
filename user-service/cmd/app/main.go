package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	http2 "user-service/internal/adapters/http"
	"user-service/internal/adapters/mysql"
	usecase "user-service/internal/app"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open(
		"mysql", "root:root@tcp(localhost:3306)/user-records")
	if db == nil || err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		panic("Failed to connect to the database")
	}
	log.Printf("Connected to the database successfully")
	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping the database: %v", err)
		panic("Failed to ping the database")
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to close the connection to the database: %v", err)
		}
	}(db)

	userRepository := mysql.NewSQLUserRepository(db)
	userInteractor := usecase.NewUserService(userRepository)
	userHandler := http2.NewUserHandler(userInteractor)

	router := mux.NewRouter()
	//router.Path()
	router.Use(loggingMiddleware)
	router.HandleFunc("/v1/users", userHandler.Create).Methods("POST")
	router.HandleFunc("/v1/users/{id}", userHandler.Get).Methods("GET")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Printf("Failed to start the server: %v", err)
		panic("Failed to start the server")
		return
	}
	fmt.Printf("Server successfully started")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: [%s] %s \"%s\"\n", time.Now().String(), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
