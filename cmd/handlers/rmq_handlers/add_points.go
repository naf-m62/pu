package rmqhandlers

import (
	"context"

	"go.uber.org/zap"

	"pu/cmd/utils"
)

type AddPointsMessage struct {
	ID     int64 `json:"id"`
	Points int32 `json:"points"`
	UserID int64 `json:"userID"`
}

func (h *Handler) AddPoints(ctx context.Context, msg []byte) (err error) {
	e := &AddPointsMessage{}
	if err = e.UnmarshalJSON(msg); err != nil {
		l := utils.GetLoggerFromContext(ctx)
		l.Error("can't unmarshal", zap.Error(err))
		return ErrInternal
	}
	if err = h.userProcessor.AddPoints(ctx, e.UserID, e.Points); err != nil {
		return GetError(err)
	}
	return nil
}
