package httpadp

import (
	"encoding/json"
	"fmt"
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
	log            log.Logger
}

func NewUserHandler(userInteractor app.UserInteractor) *UserHandler {
	return &UserHandler{UserInteractor: userInteractor, validate: NewUserValidate()}
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
		writeResponse(w, nil, shared.ServiceError{Code: shared.InvalidFormatErr, Msg: err.Error()})
		return
	}

	err := validateLoginUserRequest(h.validate, loginRequestData)
	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	loginRequestData.normalize()

	token, err := h.UserInteractor.Login(loginRequestData.Email, loginRequestData.Password)
	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	writeResponse(w, authResponse{Token: token}, nil)
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
		writeResponse(w, nil, shared.ServiceError{Code: shared.InvalidFormatErr, Msg: err.Error()})
		return
	}

	err := validateCreateUserRequest(h.validate, createUserRequestData)
	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	createUserRequestData.normalize()

	userId, err := h.UserInteractor.Create(app.RegisterUserCommand{
		Email:    createUserRequestData.Email,
		Name:     createUserRequestData.Name,
		Lastname: createUserRequestData.LastName,
		Password: createUserRequestData.Password,
	})

	writeResponse(w, userId, err)
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
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeResponse(w, nil, shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "Invalid user id in path"})
		return
	}

	getUserResult, err := h.UserInteractor.Get(id)
	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	err = json.NewEncoder(w).Encode(getUserResult)
	if err != nil {
		writeResponse(w, nil, shared.ServiceError{Code: shared.InternalErr, Msg: err.Error()})
		return
	}
}

func writeResponse(w http.ResponseWriter, v any, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		v = err
	}

	if v == nil {
		w.WriteHeader(http.StatusNoContent)
		fmt.Println("WARNING: No content to write in response")
		return
	}

	switch val := v.(type) {
	case uuid.UUID:
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(val.String())
		if err != nil {
			fmt.Println("Error encoding response body:", err)
		}

	case shared.ServiceError:
		switch val.Code {
		case shared.InvalidFormatErr, shared.InvalidRequestErr:
			w.WriteHeader(http.StatusBadRequest)
		case shared.NotFoundErr:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		err := json.NewEncoder(w).Encode(val)
		if err != nil {
			fmt.Println("Error encoding response body:", err)
		}

	case error:
		fmt.Println("Error in response:", val)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": val.Error()})

	}
}
