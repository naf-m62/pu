package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrDBRequest        = errors.New("")
	ErrWrongLoginOrPass = errors.New("")
	ErrGeneratePassword = errors.New("")
	ErrPublishEvent     = errors.New("")
	ErrDBNotFound       = errors.New("")
)

var mapError = map[error]struct{}{
	ErrDBRequest:        {},
	ErrWrongLoginOrPass: {},
	ErrGeneratePassword: {},
	ErrPublishEvent:     {},
	ErrDBNotFound:       {},
}

// GetError if err = procerror from repository, publisher ant etc. just throw it
// if err != procerror, throw procerror(cErr)
func GetError(err, cErr error) error {
	if _, ok := mapError[err]; ok {
		return err
	}
	return cErr
}
