package button

type ChangeNamePayload struct {
	DisplayName string
}

func ChangeName(caption, displayName string) (Button, error) {
	return createButton(caption, OperationChangeName,
		ChangeNamePayload{
			DisplayName: displayName,
		},
	)
}
