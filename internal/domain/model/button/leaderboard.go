package button

type LeaderboardPagePayload struct {
	PageNumber int
}

func LeaderboardPage(caption string, pageNumber int) (Button, error) {
	return CreateButton(caption, OperationLeaderboardPage, false,
		LeaderboardPagePayload{
			PageNumber: pageNumber,
		},
	)
}

func LeaderboardMyPage(caption string) Button {
	return CreateButtonWithoutPayload(caption, OperationLeaderboardMyPage, false)
}
