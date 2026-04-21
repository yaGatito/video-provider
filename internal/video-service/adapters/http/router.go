package httpadp

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	PathVarVideoID     string = "videoID"
	PathVarPublisherID string = "publisherID"

	RoutePublisherVideos string = "/v1/videos/pub/{" + PathVarPublisherID + "}"
	routeVideoSearch     string = "/v1/videos/search"
	RouteVideo           string = "/v1/videos/id/{" + PathVarVideoID + "}"
	routeSwagger         string = "/v1/swagger/"
)

func SetupRouter(
	r *mux.Router,
	h *VideoHandler,
	auth mux.MiddlewareFunc,
	logging mux.MiddlewareFunc,
	cors mux.MiddlewareFunc,
) {
	r.Use(cors)
	r.Use(logging)

	// Public routes (no auth required)
	publicRouter := r.PathPrefix("").Subrouter()

	publicRouter.HandleFunc(routeVideoSearch, h.SearchGlobal).
		Methods(http.MethodGet, http.MethodOptions)
	publicRouter.HandleFunc(RouteVideo, h.GetByID).
		Methods(http.MethodGet, http.MethodOptions)
	publicRouter.HandleFunc(RoutePublisherVideos, h.GetByPublisher).
		Methods(http.MethodGet, http.MethodOptions)

	publicRouter.PathPrefix(routeSwagger).HandlerFunc(httpSwagger.WrapHandler)

	// Protected routes (requires auth)
	protectedRouter := r.PathPrefix("").Subrouter()
	protectedRouter.Use(auth)

	protectedRouter.HandleFunc(RoutePublisherVideos, h.Create).
		Methods(http.MethodPost, http.MethodOptions)
}
