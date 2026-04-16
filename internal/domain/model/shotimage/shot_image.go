package shotimage

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

type ImageType string

const (
	ImageTypePrepareAttack = "prepare_attack"
	ImageTypePrepareDefend = "prepare_defend"
	ImageTypeAttack        = "attack"
	ImageTypeDefend        = "defend"
	ImageTypeWin           = "win"
	ImageTypeLose          = "lose"
	ImageTypeDraw          = "draw"
)

func (it ImageType) String() string {
	return string(it)
}

type ShotImage struct {
	Type                 ImageType
	AttackerExpectedSide game.ShotSide
	AttackerActualSide   game.ShotSide
	DefenderActualSide   game.ShotSide
}

func (si ShotImage) GOBEncode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(si); err != nil {
		return nil, fmt.Errorf("gob enbode: %w", err)
	}

	return buf.Bytes(), nil
}

func GOBDecode(b []byte) (ShotImage, error) {
	var si ShotImage
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&si); err != nil {
		return ShotImage{}, fmt.Errorf("gob decode: %w", err)
	}

	return si, nil
}
