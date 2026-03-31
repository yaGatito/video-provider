package httpadp

import (
	"strings"
)

type createUserRequest struct {
	Email    string `json:"email" validate:"required,email,lenLimit"`
	Name     string `json:"name" validate:"required,text64"`
	LastName string `json:"lastname" validate:"required,text64"`
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

type loginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

func (r loginUserRequest) normalize() {
	r.Email = strings.ToLower(r.Email)
}
