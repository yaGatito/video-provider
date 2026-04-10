package httpadp

import (
	"encoding/json"
	"log"
	"net/http"
	"user-service/app"
	"user-service/domain"
	"user-service/pkg/shared"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userInteractor app.UserInteractor
	validate       *validator.Validate
	log            *log.Logger
}

func NewUserHandler(userInteractor app.UserInteractor, log *log.Logger) *UserHandler {
	return &UserHandler{userInteractor: userInteractor, log: log, validate: newUserValidate()}
}

// Login godoc
//
//	@Summary		User login
//	@Tags			Users
//	@Description	Authenticate a user and return a JWT token
//	@Accept			json
//	@Produce		json
//	@Param			user	body		loginUserRequest	true	"Login user payload"
//	@Success		200		{object}	authResponse
//	@Failure		400		{object}	serviceErrorResponse
//	@Failure		401		{object}	serviceErrorResponse
//	@Failure		500		{object}	serviceErrorResponse
//	@Router			/v1/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequestData loginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequestData); err != nil {
		h.writeErrorResponse(w, shared.NewError(http.StatusBadRequest, "failed to decode login request body", err))
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

	token, err := h.userInteractor.Login(r.Context(), loginRequestData.Email, []byte(loginRequestData.Password))
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, authResponse{Token: token}, http.StatusOK)
}

// CreateUser godoc
//
//	@Summary		Creates a new user
//	@Tags			Users
//	@Description	Creates a new user and return the created user's ID. Example ID
//
//	format (UUID): 123e4567-e89b-12d3-a456-426614174000
//
//	@Accept			json
//	@Produce		json
//	@Param			user	body		createUserRequest	true	"CreateUser user payload"
//	@Success		201		{string}	string				"created user id (example: 123e4567-e89b-12d3-a456-426614174000)"
//	@Failure		400		{object}	serviceErrorResponse
//	@Failure		500		{object}	serviceErrorResponse
//	@Router			/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserRequestData createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequestData); err != nil {
		h.writeErrorResponse(w, shared.NewError(http.StatusBadRequest, "failed to decode create user request body", err))
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

	userId, err := h.userInteractor.Create(r.Context(), toDomainUser(createUserRequestData), createUserRequestData.Password)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, userId, http.StatusCreated)
}

// GetUser godoc
//
//	@Summary		Get user by ID
//	@Tags			Users
//	@Description	Retrieve user details by ID. The ID can be provided as a
//
//	UUID string (example: 123e4567-e89b-12d3-a456-426614174000) or
//	numeric identifier depending on the deployment.
//
//	@Produce		json
//	@Param 			Authorization header string true "JWT token for authentication (e.g., Bearer <token>)"
//	@Param			id	path		string	true	"User ID (example: 123e4567-e89b-12d3-a456-426614174000)"
//	@Success		200	{object}	interface{}
//	@Failure		400	{object}	serviceErrorResponse
//	@Failure		500	{object}	serviceErrorResponse
//	@Router			/v1/users/{userID} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(mux.Vars(r)[pathVarUserID])
	if err != nil {
		h.writeErrorResponse(w, shared.NewError(http.StatusBadRequest, "invalid user id format", err))
		return
	}

	if userID == uuid.Nil {
		h.writeErrorResponse(w, shared.NewError(http.StatusBadRequest, "user ID cannot be empty", nil))
		return
	}

	user, err := h.userInteractor.Get(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, toUserDto(user), http.StatusOK)
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
	case shared.Error:
		h.log.Printf("Error: %v\n", vErr)

		w.WriteHeader(int(vErr.Code))
		err := json.NewEncoder(w).Encode(serviceErrorResponse{
			Message: vErr.Message,
		})
		if err != nil {
			h.log.Println("Error encoding error response body:", err)
		}

	case validator.ValidationErrors:
		h.log.Printf("Validation request body error: %s\n", vErr[0].Error())

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(serviceErrorResponse{
			Message: "invalid field: " + vErr[0].Field(),
		})
		if err != nil {
			h.log.Println("Error validating request body:", err)
		}

	case error:
		h.log.Printf("Fallback error: %s\n", vErr.Error())

		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(serviceErrorResponse{
			Message: "video-provider error",
		})
		if err != nil {
			h.log.Println("Error encoding error response body:", err)
		}
	}
}

func toDomainUser(r createUserRequest) domain.User {
	return domain.User{
		Email:    r.Email,
		Name:     r.Name,
		LastName: r.LastName,
	}
}
