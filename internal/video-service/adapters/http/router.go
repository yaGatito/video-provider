package httpadp

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	pathVarVideoID     = "videoID"
	pathVarPublisherID = "publisherID"

	routePublisherVideos = "/v1/videos/pub/{" + pathVarPublisherID + "}"
	routeVideoSearch     = "/v1/videos/search"
	routeVideo           = "/v1/videos/id/{" + pathVarVideoID + "}"
	routeSwagger         = "/v1/swagger/"
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
	r.HandleFunc(routeVideo, h.GetByID).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc(routePublisherVideos, h.Create).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc(routePublisherVideos, h.GetByPublisher).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc(routeVideoSearch, h.SearchGlobal).
		Methods(http.MethodGet, http.MethodOptions)

	r.PathPrefix(routeSwagger).HandlerFunc(httpSwagger.WrapHandler)
}
