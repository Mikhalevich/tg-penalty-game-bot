package shotimageprovider

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

const (
	imagePathTemplate = "assets/%s.png"
)

func (sip *ShotImageProvider) Image(
	ctx context.Context,
	shot shotimage.ShotImage,
) ([]byte, error) {
	payload, err := assetsFS.ReadFile(imagePath(shot))
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return payload, nil
}

func imagePath(shot shotimage.ShotImage) string {
	return fmt.Sprintf(imagePathTemplate, imageName(shot))
}

func imageName(shot shotimage.ShotImage) string {
	switch shot.Type {
	case shotimage.ImageTypeAttackerPrepare, shotimage.ImageTypeDefenderPrepare:
		return shot.Type.String()

	case shotimage.ImageTypeShot:
		return fmt.Sprintf("%s_%s", shot.AttackerSide, shot.DefenderSide)
	}

	return ""
}
