package model

import "time"

type Platform struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
