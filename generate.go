package generate

//go:generate go tool mockgen -source=internal/domain/lbrefresher/leaderboard_refresher.go -destination=internal/domain/lbrefresher/leaderboard_refresher_mock.go -package=lbrefresher
