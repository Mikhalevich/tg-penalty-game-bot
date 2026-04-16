package shotimageprovider

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

const (
	imagePathTemplate = "assets/%s/%s.png"

	prepareFolder = "prepare"
	attackFolder  = "attack"
	defendFolder  = "defend"
	resultFolder  = "result"
	attackFile    = "attack"
	defendFile    = "defend"
	winFile       = "win"
	loseFile      = "lose"
	drawFile      = "draw"
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
	folder, fileName := imageFolderAndName(shot)

	return fmt.Sprintf(imagePathTemplate, folder, fileName)
}

// imageFodlerAndName returns image folder and filename withoud suffix.
func imageFolderAndName(shot shotimage.ShotImage) (string, string) {
	switch shot.Type {
	case shotimage.ImageTypePrepareAttack:
		return prepareFolder, attackFile

	case shotimage.ImageTypePrepareDefend:
		return prepareFolder, defendFile

	case shotimage.ImageTypeAttack:
		if shot.AttackerActualSide == game.ShotSideMiss {
			return attackFolder,
				fmt.Sprintf("%s_%s_%s",
					shot.AttackerActualSide,
					shot.AttackerExpectedSide,
					shot.DefenderActualSide,
				)
		}

		return attackFolder, fmt.Sprintf("%s_%s", shot.AttackerActualSide, shot.DefenderActualSide)

	case shotimage.ImageTypeDefend:
		if shot.AttackerActualSide == game.ShotSideMiss {
			return defendFolder,
				fmt.Sprintf("%s_%s_%s",
					shot.AttackerActualSide,
					shot.AttackerExpectedSide,
					shot.DefenderActualSide,
				)
		}

		return defendFolder, fmt.Sprintf("%s_%s", shot.AttackerActualSide, shot.DefenderActualSide)

	case shotimage.ImageTypeWin:
		return resultFolder, winFile

	case shotimage.ImageTypeLose:
		return resultFolder, loseFile

	case shotimage.ImageTypeDraw:
		return resultFolder, drawFile
	}

	return "", ""
}
