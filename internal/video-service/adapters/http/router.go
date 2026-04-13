package httpadp

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
	routeSwagger         = "/v1/swagger/"
)

func SetupRouter(
	r *mux.Router,
	h *VideoHandler,
	auth mux.MiddlewareFunc,
	logging mux.MiddlewareFunc,
	cors mux.MiddlewareFunc,
) {
	// Public routes (no auth required)
	publicRouter := r.PathPrefix("").Subrouter()
	publicRouter.Use(cors)
	publicRouter.Use(logging)

	publicRouter.HandleFunc(RouteVideoSearch, h.SearchGlobal).
		Methods(http.MethodGet, http.MethodOptions)
	publicRouter.HandleFunc(RouteVideo, h.GetByID).
		Methods(http.MethodGet, http.MethodOptions)
	publicRouter.HandleFunc(RoutePublisherVideos, h.GetByPublisher).
		Methods(http.MethodGet, http.MethodOptions)

	publicRouter.PathPrefix(routeSwagger).HandlerFunc(httpSwagger.WrapHandler)

	// Protected routes (requires auth)
	protectedRouter := r.PathPrefix("").Subrouter()
	protectedRouter.Use(cors)
	protectedRouter.Use(auth)
	protectedRouter.Use(logging)

	protectedRouter.HandleFunc(RoutePublisherVideos, h.Create).
		Methods(http.MethodPost, http.MethodOptions)
}
