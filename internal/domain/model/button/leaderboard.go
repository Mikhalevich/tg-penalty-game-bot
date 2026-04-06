package button

type LeaderboardPagePayload struct {
	PageNumber int
}

func LeaderboardPage(caption string, pageNumber int) (Button, error) {
	return CreateButton(caption, OperationLeaderboardPage,
		LeaderboardPagePayload{
			PageNumber: pageNumber,
		},
	)
}
