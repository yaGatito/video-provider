package httpadp

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	pathVarUserID = "userID"

	routeUsers   = "/v1/users"
	routeUser    = "/v1/users/{" + pathVarUserID + "}"
	routeLogin   = "/v1/users/login"
	routeSwagger = "/v1/swagger/"
)

func SetupRouter(r *mux.Router, h *UserHandler) {
	// User endpoints
	r.HandleFunc(routeUsers, h.GetUser).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc(routeUsers, h.CreateUser).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc(routeLogin, h.Login).
		Methods(http.MethodPost, http.MethodOptions)

	r.PathPrefix(routeSwagger).HandlerFunc(httpSwagger.WrapHandler)
}
