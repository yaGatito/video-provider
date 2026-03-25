package httpadp

import (
	"encoding/json"
	"log"
	"net/http"
	"video-provider/internal/pkg/shared"
	"video-provider/internal/user-service/app"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserInteractor app.UserInteractor
	validate       *validator.Validate
	log            *log.Logger
}

func NewUserHandler(userInteractor app.UserInteractor, log *log.Logger) *UserHandler {
	return &UserHandler{UserInteractor: userInteractor, log: log, validate: NewUserValidate()}
}

// Login godoc
// @Summary      User login
// @Tags         Users
// @Description  Authenticate a user and return a JWT token
// @Accept       json
// @Produce      json
// @Param        user  body    loginUserRequest  true  "Login user payload"
// @Success      200   {object}  authResponse
// @Failure      400   {object}  serviceErrorResponse
// @Failure      401   {object}  serviceErrorResponse
// @Failure      500   {object}  serviceErrorResponse
// @Router       /v1/users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequestData loginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequestData); err != nil {
		h.writeErrorResponse(w, shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "failed to decode login request body",
			Err:     err})
		return
	}

	err := h.validate.Struct(loginRequestData)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	err = validatePassword([]byte(loginRequestData.Password))
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	loginRequestData.normalize()

	token, err := h.UserInteractor.Login(loginRequestData.Email, loginRequestData.Password)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, authResponse{Token: token}, http.StatusOK)
}

// CreateUser godoc
// @Summary      CreateUser a new user
// @Tags         Users
// @Description  CreateUser a new user and return the created user's ID. Example ID
//
//	format (UUID): 123e4567-e89b-12d3-a456-426614174000
//
// @Accept       json
// @Produce      json
// @Param        user  body    createUserRequest  true  "CreateUser user payload"
// @Success      200   {string}  string  "created user id (example: 123e4567-e89b-12d3-a456-426614174000)"
// @Failure      400   {object}  serviceErrorResponse
// @Failure      500   {object}  serviceErrorResponse
// @Router       /v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserRequestData createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequestData); err != nil {
		h.writeErrorResponse(w, shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "failed to decode create user request body",
			Err:     err})
		return
	}

	err := h.validate.Struct(createUserRequestData)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	err = validatePassword([]byte(createUserRequestData.Password))
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	createUserRequestData.normalize()

	userId, err := h.UserInteractor.Create(app.RegisterUserCommand{
		Email:    createUserRequestData.Email,
		Name:     createUserRequestData.Name,
		Lastname: createUserRequestData.LastName,
		Password: createUserRequestData.Password,
	})
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, userId, http.StatusCreated)
}

// GetUser godoc
// @Summary      GetUser user by ID
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
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.writeErrorResponse(w, shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "unparsable user id"})
		return
	}

	if userID == uuid.Nil {
		h.writeErrorResponse(w, shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "empty user ID",
		})
		return
	}

	getUserResult, err := h.UserInteractor.Get(userID)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, getUserResult, http.StatusOK)
}

func (h *UserHandler) writeResponse(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		h.log.Println("Error encoding response body:", err)
	}
}

func (h *UserHandler) writeErrorResponse(w http.ResponseWriter, vErr error) {
	w.Header().Set("Content-Type", "application/json")

	switch vErr := vErr.(type) {
	case shared.ServiceError:
		h.log.Println("ServiceError: %w", vErr.Err)
		w.WriteHeader(vErr.Code)
		err := json.NewEncoder(w).Encode(serviceErrorResponse{
			Message: vErr.Message,
		})
		if err != nil {
			h.log.Println("Error encoding error response body:", err)
		}

	case validator.ValidationErrors:
		h.log.Println("Validation request body error: %w", vErr[0])
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(serviceErrorResponse{
			Message: "invalid field: " + vErr[0].Field(),
		})
		if err != nil {
			h.log.Println("Error validating request body:", err)
		}

	case error:
		h.log.Println("Fallback error: %w", vErr)
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(serviceErrorResponse{
			Message: "internal error",
		})
		if err != nil {
			h.log.Println("Error encoding error response body:", err)
		}
	}
}
