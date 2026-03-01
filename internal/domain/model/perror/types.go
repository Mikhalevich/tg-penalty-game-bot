package perror

type Type int

const (
	TypeUnspecified Type = iota
	TypeNotFound
	TypeAlreadyExists
	TypeInvalidParam
	TypeTooManyRequests
	TypeInvalidPlayer
	TypeInvalidGameState
	TypeInvalidRound
	TypeRoundNotCompleted
)
