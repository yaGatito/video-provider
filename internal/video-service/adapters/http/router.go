package httpadapter

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	PathVarVideoID     = "videoID"
	PathVarPublisherID = "publisherID"

	RoutePublisherVideos = "/v1/videos/pub/{" + PathVarPublisherID + "}"
	RouteVideoSearch     = "/v1/videos/search"
	RouteVideo           = "/v1/videos/id/{" + PathVarVideoID + "}"
)

// CORSMiddleware adds CORS headers to all responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetupRouter(r *mux.Router, h VideoHandler) {
	// Apply CORS middleware to all routes
	r.Use(CORSMiddleware)

	// Video endpoints
	r.HandleFunc(RouteVideo, h.GetByID).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc(RoutePublisherVideos, h.Create).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc(RoutePublisherVideos, h.GetByPublisher).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc(RouteVideoSearch, h.SearchGlobal).
		Methods(http.MethodGet, http.MethodOptions)

	r.PathPrefix("/v1/swagger/").HandlerFunc(httpSwagger.WrapHandler)
}
