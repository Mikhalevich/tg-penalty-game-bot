package perror

import (
	"errors"
	"time"
)

type Error struct {
	Type         Type
	Message      string
	WaitDuration time.Duration
}

func New(t Type, msg string) Error {
	return Error{
		Type:    t,
		Message: msg,
	}
}

func ParseError(err error) Error {
	var perr Error
	if errors.As(err, &perr) {
		return perr
	}

	return New(TypeUnspecified, "unspecified error")
}

func IsType(err error, t Type) bool {
	var perr Error
	if errors.As(err, &perr) {
		if perr.Type == t {
			return true
		}
	}

	return false
}

func IsTypeWithError(err error, t Type) (Error, bool) {
	var perr Error
	if errors.As(err, &perr) {
		if perr.Type == t {
			return perr, true
		}
	}

	return Error{}, false
}

func (e Error) Error() string {
	return e.Message
}

func NotFound(msg string) Error {
	return New(TypeNotFound, msg)
}

func AlreadyExists(msg string) Error {
	return New(TypeAlreadyExists, msg)
}

func InvalidParam(msg string) Error {
	return New(TypeInvalidParam, msg)
}

func TooManyRequests(msg string, waitDuration time.Duration) Error {
	return Error{
		Type:         TypeTooManyRequests,
		Message:      msg,
		WaitDuration: waitDuration,
	}
}

func InvalidPlayer() Error {
	return New(TypeInvalidPlayer, "invalid player")
}
