package httpadp

import (
	"encoding/json"
	"errors"
	"examples-user-service/internal/app"
	"examples-user-service/internal/shared"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserInteractor app.UserServiceImpl
}

func NewUserHandler(userInteractor app.UserServiceImpl) UserHandler {
	return UserHandler{UserInteractor: userInteractor}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUserRequestData createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequestData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		writeJSON(w, http.StatusBadRequest, serviceErrorResponse{Code: shared.ServiceErrorCodeInvalidRequest})
		return
	}

	if err := createUserRequestData.validate(); err != nil {
		log.Printf("Validation error: %v", err)
		var vErr shared.ValidationError
		errors.As(err, &vErr)
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
