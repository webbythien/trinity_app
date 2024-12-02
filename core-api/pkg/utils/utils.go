package utils

import (
	"strings"

	"github.com/google/uuid"
)

func NewID() string {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")
	return id
}
