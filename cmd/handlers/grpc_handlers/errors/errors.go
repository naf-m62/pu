package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	procerrors "pu/cmd/processor/errors"
)

var (
	ErrInternal         = status.Error(codes.Internal, "internal error")
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "wrong login or password")
)

var mapProcError = map[error]error{
	procerrors.ErrDBRequest:        ErrInternal,
	procerrors.ErrWrongLoginOrPass: ErrPermissionDenied,
	procerrors.ErrGeneratePassword: ErrInternal,
	procerrors.ErrPublishEvent:     ErrInternal,
}

func GetError(err error) error {
	if _, ok := mapProcError[err]; !ok {
		return ErrInternal
	}
	return mapProcError[err]
}
