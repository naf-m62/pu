package httphandlers

import (
	"net/http"

	"github.com/pkg/errors"

	procerrors "pu/cmd/processor/errors"
)

type Error struct {
	Err  error
	Code int
}

var (
	ErrInternal = Error{
		Err:  errors.New("internal error"),
		Code: http.StatusInternalServerError,
	}
	ErrBadRequest = Error{
		Err:  errors.New("wrong request"),
		Code: http.StatusBadRequest,
	}
)

var mapProcError = map[error]Error{
	procerrors.ErrDBRequest:    ErrInternal,
	procerrors.ErrPublishEvent: ErrInternal,
}

func GetError(err error) Error {
	if _, ok := mapProcError[err]; !ok {
		return ErrInternal
	}
	return mapProcError[err]
}
