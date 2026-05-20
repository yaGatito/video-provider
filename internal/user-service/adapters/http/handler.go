package httpadp

import (
	"encoding/json"
	"net/http"
	"video-provider/pkg/common"
	"video-provider/user-service/app"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userInteractor app.UserInteractor
	validate       *validator.Validate
	log            *common.Logger
}

func NewUserHandler(userInteractor app.UserInteractor, log *common.Logger) *UserHandler {
	validate, err := newUserValidate()
	if err != nil {
		log.Error("validate is not created. aborting..", err)
	}
	return &UserHandler{userInteractor: userInteractor, log: log, validate: validate}
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
		common.WriteErrorResponse(w, h.log, &common.Error{
			Code: http.StatusBadRequest, Message: "failed to decode login request body", Err: err,
		})
		return
	}

	err := h.validate.Struct(loginRequestData)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	err = validatePassword([]byte(loginRequestData.Password))
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	loginRequestData.normalize()

	token, err := h.userInteractor.Login(
		r.Context(),
		loginRequestData.Email,
		[]byte(loginRequestData.Password),
	)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	common.WriteResponse(w, h.log, authResponse{Token: token}, http.StatusOK)
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
		common.WriteErrorResponse(w, h.log, &common.Error{
			Err:     err,
			Code:    http.StatusBadRequest,
			Message: "failed to decode create user request body",
		})
		return
	}

	err := h.validate.Struct(createUserRequestData)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	err = validatePassword([]byte(createUserRequestData.Password))
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	createUserRequestData.normalize()

	userId, err := h.userInteractor.Create(
		r.Context(),
		toDomainUser(createUserRequestData),
		createUserRequestData.Password,
	)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	common.WriteResponse(w, h.log, createUserResponse{
		UserID: userId.String(),
	}, http.StatusCreated)
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
		common.WriteErrorResponse(w, h.log, &common.Error{
			Err:     err,
			Code:    http.StatusBadRequest,
			Message: "invalid user id format",
		})
		return
	}

	if userID == uuid.Nil {
		common.WriteErrorResponse(w, h.log, &common.Error{
			Code:    http.StatusBadRequest,
			Message: "user ID cannot be empty",
		})
		return
	}

	user, err := h.userInteractor.Get(r.Context(), userID)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	common.WriteResponse(w, h.log, toUserDto(user), http.StatusOK)
}
