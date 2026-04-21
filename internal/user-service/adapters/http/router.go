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
	routeLogin   = "/v1/login"
	routeSwagger = "/v1/swagger/"
)

func SetupRouter(
	r *mux.Router,
	h *UserHandler,
	auth mux.MiddlewareFunc,
	logging mux.MiddlewareFunc,
	cors mux.MiddlewareFunc,
) {
	r.Use(cors)
	r.Use(logging)

	// Public routes (no auth required)
	publicRouter := r.PathPrefix("").Subrouter()

	// Login and register routes (public)
	publicRouter.HandleFunc(routeLogin, h.Login).
		Methods(http.MethodPost, http.MethodOptions)
	publicRouter.HandleFunc(routeUsers, h.CreateUser).
		Methods(http.MethodPost, http.MethodOptions)

	// Swagger route (public)
	publicRouter.PathPrefix(routeSwagger).HandlerFunc(httpSwagger.WrapHandler)

	// Protected routes (requires auth)
	protectedRouter := r.PathPrefix("").Subrouter()
	protectedRouter.Use(auth)

	// User endpoints (protected)
	protectedRouter.HandleFunc(routeUser, h.GetUser).
		Methods(http.MethodGet, http.MethodOptions)
}
