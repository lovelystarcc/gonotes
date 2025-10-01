package dto

import (
	"fmt"
	"gonotes/internal/auth/entity"
	"net/http"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserRequest) Bind(r *http.Request) error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if u.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		Email: user.Email,
		ID:    user.ID,
	}
}

func (n *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (n *LoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewLoginResponse(email, token string) *LoginResponse {
	return &LoginResponse{Email: email, Token: token}
}
