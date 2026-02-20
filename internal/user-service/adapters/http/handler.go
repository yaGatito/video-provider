package httpadp

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"video-provider/internal/user-service/app"
	"video-provider/internal/user-service/shared"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserInteractor *app.UserService
	log            log.Logger
}

func NewUserHandler(userInteractor *app.UserService) UserHandler {
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
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := createUserRequestData.validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
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
		var vErr shared.ServiceError
		if errors.As(err, &vErr) {
			log.Printf("Error registering user (validation): %v", err)
			writeJSON(w, http.StatusBadRequest, serviceErrorResponse{Code: shared.ValidationErr, Payload: vErr})
			return
		}
		log.Printf("Error registering user: %v", err)
		writeJSON(w, http.StatusInternalServerError, serviceErrorResponse{Code: shared.InternalErr})
		return
	}

	err = json.NewEncoder(w).Encode(userId)
	if err != nil {
		log.Printf("Error encoding user response: %v", err)
		writeJSON(w, http.StatusInternalServerError, serviceErrorResponse{Code: shared.InternalErr})
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
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, serviceErrorResponse{Code: shared.InvalidRequestErr})
		log.Printf("Invalid user id in path: %v", err)
		return
	}
	log.Printf("Find by id: %s", id.String())

	getUserResult, err := h.UserInteractor.Get(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, serviceErrorResponse{Code: shared.NotFoundErr})
		log.Printf("Error retrieving user: %v", err)
		return
	}
	log.Printf("User by id: %s found! - name: %s, lastname: %s\n", id.String(), getUserResult.Name, getUserResult.Lastname)

	err = json.NewEncoder(w).Encode(getUserResult)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, serviceErrorResponse{Code: shared.InternalErr})
		log.Printf("Error encoding user response: %v", err)
		return
	}
	log.Println("Response were written successfully")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	var vErr shared.ServiceError
	resp := serviceErrorResponse{}
	if errors.As(err, &vErr) {
		resp.Code = vErr.Code
		resp.Payload = vErr.Msg
	} else {
		resp.Code = shared.InternalErr
		// Verbose way. TODO: change next on complete
		resp.Payload = err.Error()
		// resp.Payload = "internal error"
	}
	log.Printf("Error writing response: %v", err)
	writeJSON(w, status, resp)
}
