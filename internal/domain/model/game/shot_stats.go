package game

type ShotStats struct {
	Attack []ShotSidePercent
	Defend []ShotSidePercent
}

type ShotSidePercent struct {
	Side    ShotSide
	Percent float32
}
