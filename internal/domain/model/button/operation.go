package button

type Operation string

const (
	OperationChangeName      Operation = "ChangeName"
	OperationShotSide        Operation = "ShotSide"
	OperationLeaveGame       Operation = "LeaveGame"
	OperationStopSearchGame  Operation = "StopSearchGame"
	OperationLeaderboardPage Operation = "LeaderboardPage"
)
