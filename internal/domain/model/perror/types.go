package perror

type Type int

const (
	TypeUnspecified Type = iota
	TypeNotFound
	TypeAlreadyExists
	TypeNotExists
	TypeInvalidParam
	TypeTooManyRequests
	TypeInvalidPlayer
	TypeInvalidGameState
	TypeInvalidRound
	TypeRoundNotCompleted
	TypeAlreadyInGame
	TypeNotInGame
	TypeInSearchGameState
)
