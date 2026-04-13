package button

type ChangeNamePayload struct {
	DisplayName string
}

func ChangeName(caption, displayName string) (Button, error) {
	return CreateButton(caption, OperationChangeName,
		ChangeNamePayload{
			DisplayName: displayName,
		},
	)
}

func ChangeNameTrigger(caption string) Button {
	return CreateButtonWithoutPayload(caption, OperationChangeNameTrigger)
}
