package button

type Operation string

const (
	OperationChangeName        Operation = "ChangeName"
	OperationChangeNameTrigger Operation = "ChangeNameTrigger"
	OperationChangeNameCancel  Operation = "ChangeNameCancel"
	OperationShotSide          Operation = "ShotSide"
	OperationLeaveGame         Operation = "LeaveGame"
	OperationShotStats         Operation = "ShotStats"
	OperationStopSearchGame    Operation = "StopSearchGame"
	OperationLeaderboardPage   Operation = "LeaderboardPage"
	OperationRepeatGame        Operation = "RepeatGame"
)

func (o Operation) String() string {
	return string(o)
}
