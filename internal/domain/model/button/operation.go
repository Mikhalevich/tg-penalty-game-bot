package button

type Operation string

const (
	OperationChangeName      Operation = "ChangeName"
	OperationShotSide        Operation = "ShotSide"
	OperationLeaveGame       Operation = "LeaveGame"
	OperationShotStats       Operation = "ShotStats"
	OperationStopSearchGame  Operation = "StopSearchGame"
	OperationLeaderboardPage Operation = "LeaderboardPage"
)
