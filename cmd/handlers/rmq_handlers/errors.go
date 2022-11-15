package rmqhandlers

import (
	"github.com/pkg/errors"

	procerrors "pu/cmd/processor/errors"
)

var (
	ErrInternal = errors.New("internal error")
)

var mapProcError = map[error]error{
	procerrors.ErrDBRequest: ErrInternal,
}

func GetError(err error) error {
	if _, ok := mapProcError[err]; !ok {
		return ErrInternal
	}
	return mapProcError[err]
}
