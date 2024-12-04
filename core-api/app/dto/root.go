package dto

import (
	"time"

	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
)

// ###########################
// ## Data Transfer Objects ##
// ###########################
type Response struct {
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	Success   bool        `json:"success"`
	ErrorCode *string     `json:"error_code"`
}

type User struct {
	ID        int        `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	IsActive  bool       `json:"is_active"`
	IsAdmin   bool       `json:"is_admin"`
	UserName  string     `json:"username"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
}

func ToUser(u *model.User) *User {
	return &User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsActive:  u.IsActive,
		IsAdmin:   u.IsAdmin,
		UserName:  u.UserName,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func ToUsers(users []*model.User) []*User {
	res := make([]*User, len(users))
	for i, user := range users {
		res[i] = ToUser(user)
	}
	return res
}
