package user

import (
	"context"

	"go.uber.org/zap"

	procerrors "pu/cmd/processor/errors"
	"pu/cmd/utils"
)

func (p *Processor) AddPoints(ctx context.Context, userID int64, points int32) (err error) {
	l := utils.GetLoggerFromContext(ctx)

	if err = p.userRepo.AddPoints(ctx, userID, points); err != nil {
		l.Error("can't add point", zap.Error(err))
		return procerrors.GetError(err, procerrors.ErrDBRequest)
	}
	l.Info("points added to userID", zap.Int64("userID", userID))
	return nil
}
