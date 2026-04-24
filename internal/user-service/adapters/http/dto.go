package httpadp

import (
	"strings"
	"time"
	"video-provider/user-service/domain"
)

type createUserRequest struct {
	Email    string `json:"email"    validate:"required,email,lenLimit"`
	Name     string `json:"name"     validate:"required,text64"`
	LastName string `json:"lastname" validate:"required,text64"`
	Password string `json:"password" validate:"required"`
}

type createUserResponse struct {
	UserID string `json:"user_id"`
}

type loginUserRequest struct {
	Email    string `json:"email"    validate:"required,email,lenLimit"`
	Password string `json:"password" validate:"required"`
}

func (r createUserRequest) normalize() {
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	r.Name = strings.TrimSpace(r.Name)
	r.LastName = strings.TrimSpace(r.LastName)
}

type serviceErrorResponse struct {
	Message string `json:"msg"`
}

type authResponse struct {
	Token string `json:"token"`
}

type userResponse struct {
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func (r loginUserRequest) normalize() {
	r.Email = strings.ToLower(r.Email)
}

func toUserDto(u domain.User) userResponse {
	return userResponse{
		Name:      u.Name,
		Lastname:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
