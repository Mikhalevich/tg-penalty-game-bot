package game

type GameStatus string

const (
	GameStatusInProgress GameStatus = "in_progress"
	GameStatusCompleted  GameStatus = "completed"
)

func (gs GameStatus) String() string {
	return string(gs)
}
