package main

import (
	"log"
	"net/http"
	"os"
	"time"
	http2 "video-service/internal/adapters/http"
	"video-service/internal/adapters/idgen"
	"video-service/internal/adapters/testdb"
	"video-service/internal/app"
	"video-service/internal/domain"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// ctx := context.Background()
	// conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
	// if err != nil {
	// 	return err
	// }
	// defer conn.Close(ctx)

	// videoRepository := postgres.NewVideoRepoPostgreSQL(conn)

	idGen := idgen.New()
	mwLog := MiddlewareLogger{log: log.New(os.Stdout, "[VSRVC] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)}
	store := make(map[uuid.UUID]domain.Video)

	videoRepository := testdb.NewVideoRepoTestDB(store, mwLog.Log())
	videoService := app.NewVideoInteractor(videoRepository)
	videoHandler := http2.NewVideoHandler(videoService, idGen, mwLog.log)

	router := mux.NewRouter()
	router.Use(mwLog.loggingMiddleware)
	router.HandleFunc("/v1/videos/{videoId}", videoHandler.GetByID).Methods("GET")
	router.HandleFunc("/v1/videos/pub/{publisherId}", videoHandler.Create).Methods("POST")
	router.HandleFunc("/v1/videos/pub/{publisherId}", videoHandler.GetByPublisher).Methods("GET")
	router.HandleFunc("/v1/videos/search/", videoHandler.SearchGlobal).Methods("GET")

	mwLog.log.Printf("Server successfully started")
	err := http.ListenAndServe(":8081", router)
	if err != nil {
		return err
	}
	return nil
}

type MiddlewareLogger struct {
	log *log.Logger
}

func (l *MiddlewareLogger) Log() *log.Logger {
	return l.log
}

func (l *MiddlewareLogger) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: [%s] %s \"%s\"\n", time.Now().String(), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
