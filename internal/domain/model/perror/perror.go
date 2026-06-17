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

func ParseError(err error) (Error, bool) {
	var perr Error
	if errors.As(err, &perr) {
		return perr, true
	}

	return New(TypeUnspecified, "unspecified error"), false
}

func IsType(err error, t Type) bool {
	if perr, ok := errors.AsType[Error](err); ok {
		if perr.Type == t {
			return true
		}
	}

	return false
}

func IsTypeWithError(err error, t Type) (Error, bool) {
	if perr, ok := errors.AsType[Error](err); ok {
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

func NotExists(msg string) Error {
	return New(TypeNotExists, msg)
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

func InvalidState(msg string) Error {
	return New(TypeInvalidState, msg)
}

func InvalidGameState() Error {
	return New(TypeInvalidGameState, "invalid game state")
}

func InvalidRound() Error {
	return New(TypeInvalidRound, "invalid round")
}

func RoundNotCompleted() Error {
	return New(TypeRoundNotCompleted, "round not completed")
}

func AlreadyInGame() Error {
	return New(TypeAlreadyInGame, "already in game")
}

func NotInGame() Error {
	return New(TypeNotInGame, "not in game")
}

func InSearchGameState() Error {
	return New(TypeInSearchGameState, "searching for game")
}

func NoRowsUpdated() Error {
	return New(TypeNoRowsUpdated, "no rows updated")
}
