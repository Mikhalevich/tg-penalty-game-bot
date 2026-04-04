package button

func LeaderboardPlayer(caption string) Button {
	return CreateButtonWithoutPayload(caption, OperationLeaderboardPlayer)
}

func LeaderboardTop(caption string) Button {
	return CreateButtonWithoutPayload(caption, OperationLeaderboardTop)
}
