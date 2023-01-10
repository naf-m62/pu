package database

import (
	"context"

	"github.com/jackc/pgx/v5"

	"pu/postgres"
)

type TxRepo struct {
	db postgres.DB
}

func NewTx(db postgres.DB) *TxRepo {
	return &TxRepo{db: db}
}

func (t *TxRepo) Do(ctx context.Context, fn func(tx pgx.Tx) error) (err error) {
	var tx pgx.Tx
	if tx, err = t.db.DB().Begin(ctx); err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		if errR := tx.Rollback(ctx); errR != nil {
			return err
		}
		return err
	}
	if errC := tx.Commit(ctx); errC != nil {
		return errC
	}
	return nil
}
