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

func SetupRouter(r *mux.Router, h *VideoHandler, auth mux.MiddlewareFunc, logging mux.MiddlewareFunc, cors mux.MiddlewareFunc) {
	// Public routes (no auth required)
	publicRouter := r.PathPrefix("").Subrouter()
	publicRouter.Use(cors)
	publicRouter.Use(logging)

	// Global public search
	publicRouter.HandleFunc(routeVideoSearch, h.SearchGlobal).
		Methods(http.MethodGet, http.MethodOptions)

	// Swagger route (public)
	publicRouter.PathPrefix(routeSwagger).HandlerFunc(httpSwagger.WrapHandler)

	// Protected routes (requires auth)
	protectedRouter := r.PathPrefix("").Subrouter()
	protectedRouter.Use(cors)
	protectedRouter.Use(auth)
	protectedRouter.Use(logging)

	// Video endpoints (protected)
	protectedRouter.HandleFunc(routeVideo, h.GetByID).
		Methods(http.MethodGet, http.MethodOptions)
	protectedRouter.HandleFunc(routePublisherVideos, h.Create).
		Methods(http.MethodPost, http.MethodOptions)
	protectedRouter.HandleFunc(routePublisherVideos, h.GetByPublisher).
		Methods(http.MethodGet, http.MethodOptions)
}
