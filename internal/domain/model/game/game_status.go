package game

type GameStatus string

const (
	GameStatusInProgress GameStatus = "in_progress"
	GameStatusFinished   GameStatus = "finished"
	GameStatusCanceled   GameStatus = "canceled"
)

func (gs GameStatus) String() string {
	return string(gs)
}
