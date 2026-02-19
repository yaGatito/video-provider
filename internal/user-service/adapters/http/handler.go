package httpadp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"video-provider/internal/user-service/app"
	"video-provider/internal/user-service/shared"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserInteractor app.UserService
}

func NewUserHandler(userInteractor app.UserService) UserHandler {
	return UserHandler{UserInteractor: userInteractor}
}

// Create godoc
// @Summary      Create a new user
// @Tags         Users
// @Description  Create a new user and return the created user's ID. Example ID
//
//	format (UUID): 123e4567-e89b-12d3-a456-426614174000
//
// @Accept       json
// @Produce      json
// @Param        user  body    createUserRequest  true  "Create user payload"
// @Success      200   {string}  string  "created user id (example: 123e4567-e89b-12d3-a456-426614174000)"
// @Failure      400   {object}  serviceErrorResponse
// @Failure      500   {object}  serviceErrorResponse
// @Router       /v1/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUserRequestData createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequestData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		writeJSON(w, http.StatusBadRequest, serviceErrorResponse{Code: shared.ServiceErrorCodeInvalidRequest})
		return
	}

	var vErr shared.ValidationError
	if err := createUserRequestData.validate(); err != nil && errors.As(err, &vErr) {
		log.Printf("Validation error: %v", err)
		writeJSON(w, http.StatusBadRequest, serviceErrorResponse{Code: vErr.Error(), Payload: vErr.Violations})
		return
	}

	createUserRequestData.normalize()

	userId, err := h.UserInteractor.Register(app.RegisterUserCommand{
		Email:    createUserRequestData.Email,
		Name:     createUserRequestData.Name,
		Lastname: createUserRequestData.LastName,
		Password: createUserRequestData.Password,
	})

	if err != nil {
		log.Printf("Error registering user: %v", err)
		var valError shared.ValidationError
		if errors.As(err, &valError) {
			log.Printf("Validation error: %v", valError)
			writeJSON(w, http.StatusBadRequest, serviceErrorResponse{valError.Error(), valError.Violations})
			return
		}
	}

	err = json.NewEncoder(w).Encode(userId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, serviceErrorResponse{Code: shared.ServiceErrorCodeInternalError})
		log.Printf("Error encoding user response: %v", err)
		return
	}
}

// Get godoc
// @Summary      Get user by ID
// @Tags         Users
// @Description  Retrieve user details by ID. The ID can be provided as a
//
//	UUID string (example: 123e4567-e89b-12d3-a456-426614174000) or
//	numeric identifier depending on the deployment.
//
// @Produce      json
// @Param        id   path    string  true  "User ID (example: 123e4567-e89b-12d3-a456-426614174000)"
// @Success      200  {object}  interface{}
// @Failure      400  {object}  serviceErrorResponse
// @Failure      500  {object}  serviceErrorResponse
// @Router       /v1/users/{id} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	log.Printf("Find by id: %d", id)

	getUserResult, err := h.UserInteractor.Get(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("User by id: %d found! - name: %s, lastname: %s\n", id, getUserResult.Name, getUserResult.Lastname)

	err = json.NewEncoder(w).Encode(getUserResult)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding user response: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("Response were written successfully")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
