package model

import "github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"

type Button struct {
	ID              string `db:"id"`
	Caption         string `db:"caption"`
	Operation       string `db:"operation"`
	IsDeleteMessage bool   `db:"is_delete_message"`
	Payload         []byte `db:"payload"`
}

func (b *Button) ToDomButton() *button.Button {
	return &button.Button{
		ID:                   button.IDFromString(b.ID),
		Caption:              b.Caption,
		Operation:            button.Operation(b.Operation),
		IsDeleteAfterProcess: b.IsDeleteMessage,
		Payload:              b.Payload,
	}
}

func toDBButton(domButton button.Button) Button {
	return Button{
		ID:              domButton.ID.String(),
		Caption:         domButton.Caption,
		Operation:       domButton.Operation.String(),
		IsDeleteMessage: domButton.IsDeleteAfterProcess,
		Payload:         domButton.Payload,
	}
}

func ToDBButtons(rows []button.ButtonRow) []Button {
	capacity := 0
	for _, row := range rows {
		capacity += len(row)
	}

	buttons := make([]Button, 0, capacity)

	for _, row := range rows {
		for _, btn := range row {
			buttons = append(buttons, toDBButton(btn))
		}
	}

	return buttons
}
