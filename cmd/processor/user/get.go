package user

import (
	"context"

	"go.uber.org/zap"

	procerrors "pu/cmd/processor/errors"
	"pu/cmd/processor/user/domain"
	"pu/cmd/utils"
)

func (p *Processor) Get(ctx context.Context, id int64) (ru *domain.User, err error) {
	l := utils.GetLoggerFromContext(ctx)

	if ru, err = p.userRepo.Get(ctx, id); err != nil {
		l.Error("can't get user", zap.Error(err))
		return nil, procerrors.GetError(err, procerrors.ErrDBRequest)
	}
	return ru, nil
}
