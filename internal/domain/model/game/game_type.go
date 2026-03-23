package game

type GameType string

const (
	GameTypeRating   GameType = "rating"
	GameTypeFriendly GameType = "friendly"
)

func (gt GameType) String() string {
	return string(gt)
}
