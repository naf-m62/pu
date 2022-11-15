package user

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	procerrors "pu/cmd/processor/errors"
	"pu/cmd/processor/user/domain"
	"pu/cmd/utils"
)

func (p *Processor) Auth(ctx context.Context, email, password string) (ru *domain.User, err error) {
	l := utils.GetLoggerFromContext(ctx)

	if ru, err = p.userRepo.GetByEmail(ctx, email); err != nil {
		l.Error("can't get user by email", zap.Error(err))
		return nil, procerrors.ErrDBRequest
	}

	if err = bcrypt.CompareHashAndPassword([]byte(ru.PasswordHash), []byte(ru.Salt+password+saltCode)); err != nil {
		l.Error("incorrect login or password", zap.Error(err))
		return nil, procerrors.ErrWrongLoginOrPass
	}
	return ru, nil
}
