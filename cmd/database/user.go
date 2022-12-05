package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"pu/cmd/database/entities"
	procerrors "pu/cmd/processor/errors"
	"pu/cmd/processor/user/domain"
	"pu/postgres"
)

type UserRepo struct {
	db postgres.DB
}

func NewUser(db postgres.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) Create(ctx context.Context, user *domain.User) (ru *domain.User, err error) {
	iu := entities.ConvertFromUserProcessor(user)
	now := time.Now()
	uDB := &entities.User{
		Name:         iu.Name,
		Email:        iu.Email,
		CreatedAt:    now,
		UpdatedAt:    now,
		PasswordHash: iu.PasswordHash,
		Salt:         iu.Salt,
	}
	err = u.db.DB().QueryRow(ctx, `
INSERT INTO users(name, email, created_at, updated_at, password_hash, salt) VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING id
`, uDB.Name, uDB.Email, uDB.CreatedAt, uDB.UpdatedAt, uDB.PasswordHash, uDB.Salt).
		Scan(&uDB.ID)
	return uDB.ConvertToUserProcessor(), err
}

func (u *UserRepo) Get(ctx context.Context, userId int64) (user *domain.User, err error) {
	uDB := &entities.User{}
	err = u.db.DB().QueryRow(ctx, `
SELECT id, name, email, points, password_hash, salt FROM users WHERE id = $1
`, userId).Scan(&uDB.ID, &uDB.Name, &uDB.Email, &uDB.Points, &uDB.PasswordHash, &uDB.Salt)
	if err == pgx.ErrNoRows {
		return nil, procerrors.ErrDBNotFound
	}
	return uDB.ConvertToUserProcessor(), err
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (user *domain.User, err error) {
	uDB := &entities.User{}
	err = u.db.DB().QueryRow(ctx, `
SELECT id, name, email, points, password_hash, salt FROM users WHERE email = $1
`, email).Scan(&uDB.ID, &uDB.Name, &uDB.Email, &uDB.Points, &uDB.PasswordHash, &uDB.Salt)
	if err == pgx.ErrNoRows {
		return nil, procerrors.ErrDBNotFound
	}
	return uDB.ConvertToUserProcessor(), err
}

func (u *UserRepo) AddPoints(ctx context.Context, userID int64, points int32) (err error) {
	_, err = u.db.DB().Exec(ctx, `
UPDATE users SET points = points + $1, updated_at = $2 WHERE id = $3
`, points, time.Now(), userID)
	return err
}
