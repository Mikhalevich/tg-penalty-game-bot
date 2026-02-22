package player

type GameStatus string

const (
	GameStatusIdle         GameStatus = "idle"
	GameStatusReadyForGame GameStatus = "ready_for_game"
	GameStatusInGame       GameStatus = "in_game"
)

func (gs GameStatus) String() string {
	return string(gs)
}
