package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrDBRequest        = errors.New("")
	ErrWrongLoginOrPass = errors.New("")
	ErrGeneratePassword = errors.New("")
	ErrPublishEvent     = errors.New("")
)
