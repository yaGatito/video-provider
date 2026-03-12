package httpadp

import (
	"net/http"

	"video-provider/internal/pkg/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	PathVarUserID = "userID"

	RouteUsers = "/v1/users"
	RouteUser  = "/v1/users/{" + PathVarUserID + "}"
	RouteLogin = "/v1/users/login"
)

func SetupRouter(r *mux.Router, h *UserHandler) {
	// Apply CORS middleware to all routes
	r.Use(middleware.CORSMiddleware)

	// User endpoints
	r.HandleFunc(RouteUsers, h.GetUser).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc(RouteUsers, h.CreateUser).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc(RouteLogin, h.Login).
		Methods(http.MethodPost, http.MethodOptions)

	r.PathPrefix("/v1/swagger/").HandlerFunc(httpSwagger.WrapHandler)
}
