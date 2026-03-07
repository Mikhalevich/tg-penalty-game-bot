package shotimage

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

type ImageType string

const (
	ImageTypeAttackerPrepare = "attacker_prepare"
	ImageTypeDefenderPrepare = "defender_prepare"
	ImageTypeShot            = "shot"
)

func (it ImageType) String() string {
	return string(it)
}

type ShotImage struct {
	Type         ImageType
	AttackerSide game.ShotSide
	DefenderSide game.ShotSide
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
