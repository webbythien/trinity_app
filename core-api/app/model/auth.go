package model

import "time"

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
	FullName string `json:"full_name" validate:"required" example:"Nguyen Van A"`
	RoleID   int    `json:"-"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email" example:"admin@admin.com"`
	Password string `json:"password" validate:"required" example:"admin"`
}

type AuthOutput struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	RoleID    int       `json:"role_id"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
	User  *AuthOutput `json:"user"`
	Token string      `json:"access_token"`
}
