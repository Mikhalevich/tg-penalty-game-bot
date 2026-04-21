package button

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/google/uuid"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(s string) ID {
	return ID(s)
}

type Button struct {
	ID                   ID
	Caption              string
	Operation            Operation
	IsDeleteAfterProcess bool
	Payload              []byte
}

type ButtonRow []Button

func Row(buttons ...Button) ButtonRow {
	return buttons
}

func GetPayload[P any](b Button) (P, error) {
	payload, err := gobDecodePayload[P](b.Payload)
	if err != nil {
		return payload, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}

func gobEncodePayload(p any) ([]byte, error) {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}

	return buf.Bytes(), nil
}

func gobDecodePayload[Payload any](b []byte) (Payload, error) {
	var payload Payload
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&payload); err != nil {
		return payload, fmt.Errorf("gob decode: %w", err)
	}

	return payload, nil
}

func CreateButton[P any](
	caption string,
	operation Operation,
	isDelete bool,
	payload P,
) (Button, error) {
	payloadBytes, err := gobEncodePayload(payload)
	if err != nil {
		return Button{}, fmt.Errorf("encode payload: %w", err)
	}

	return Button{
		ID:                   IDFromString(generateID()),
		Caption:              caption,
		Operation:            operation,
		IsDeleteAfterProcess: isDelete,
		Payload:              payloadBytes,
	}, nil
}

func CreateButtonWithoutPayload(
	caption string,
	operation Operation,
	isDelete bool,
) Button {
	return Button{
		ID:                   IDFromString(generateID()),
		Caption:              caption,
		Operation:            operation,
		IsDeleteAfterProcess: isDelete,
	}
}

func generateID() string {
	return uuid.NewString()
}
