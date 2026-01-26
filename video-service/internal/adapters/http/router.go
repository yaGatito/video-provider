package httpadapter

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	RouteVideo           = "/v1/videos/{" + PathVarVideoID + "}"
	RoutePublisherVideos = "/v1/videos/pub/{" + PathVarPublisherID + "}"
	RouteVideoSearch     = "/v1/videos/search/"
)

func SetupRouter(r *mux.Router, h VideoHandler) {
	r.HandleFunc(RouteVideo, h.GetByID).
		Methods(http.MethodGet)

	r.HandleFunc(RoutePublisherVideos, h.Create).
		Methods(http.MethodPost)

	r.HandleFunc(RoutePublisherVideos, h.GetByPublisher).
		Methods(http.MethodGet)

	r.HandleFunc(RouteVideoSearch, h.SearchGlobal).
		Methods(http.MethodGet)

	r.PathPrefix("/swagger/").HandlerFunc(httpSwagger.WrapHandler)
	// r.Get("/swagger/*any").HandlerFunc(ginSwagger.WrapHandler(swaggerFiles.Handler))
}
