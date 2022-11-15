package domain

import (
	"time"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Points       int64
	PasswordHash string
	Salt         string
}
