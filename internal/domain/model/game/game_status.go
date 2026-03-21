package game

type GameStatus string

const (
	GameStatusPending    GameStatus = "pending"
	GameStatusInProgress GameStatus = "in_progress"
	GameStatusCompleted  GameStatus = "completed"
	GameStatusCanceled   GameStatus = "canceled"
)

func (gs GameStatus) String() string {
	return string(gs)
}
