package entities

import (
	"database/sql"
	"time"

	"pu/cmd/processor/user/domain"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Points       sql.NullInt64
	PasswordHash string
	Salt         string
}

func (u *User) ConvertToUserProcessor() *domain.User {
	if u == nil {
		return &domain.User{}
	}
	return &domain.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Points:       u.Points.Int64,
		PasswordHash: u.PasswordHash,
		Salt:         u.Salt,
	}
}

func ConvertFromUserProcessor(up *domain.User) *User {
	return &User{
		ID:           up.ID,
		Name:         up.Name,
		Email:        up.Email,
		CreatedAt:    up.CreatedAt,
		UpdatedAt:    up.UpdatedAt,
		Points:       sql.NullInt64{Valid: up.Points > 0, Int64: up.Points},
		PasswordHash: up.PasswordHash,
		Salt:         up.Salt,
	}
}
