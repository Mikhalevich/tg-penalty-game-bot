package game

type GameType string

const (
	GameTypeRating GameType = "rating"
	GameTypeByLink GameType = "by_link"
	GameTypeBot    GameType = "bot"
)

func (gt GameType) String() string {
	return string(gt)
}
