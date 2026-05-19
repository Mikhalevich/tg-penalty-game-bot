package shotimageprovider

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

const (
	imageFolderTemplate = "assets/%s/%s"

	prepareFolder = "prepare"
	attackFolder  = "attack"
	defendFolder  = "defend"
	resultFolder  = "result"
	winFolder     = "win"
	loseFolder    = "lose"
	drawFolder    = "draw"
)

func (sip *ShotImageProvider) Image(
	ctx context.Context,
	shot shotimage.ShotImage,
) ([]byte, error) {
	imagePath, err := imageAbsPath(shot)
	if err != nil {
		return nil, fmt.Errorf("iamge abs path: %w", err)
	}
	payload, err := assetsFS.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return payload, nil
}

func imageAbsPath(shot shotimage.ShotImage) (string, error) {
	var (
		folder, subFolder = imageFolderAndSubfolder(shot)
		folderAbsPath     = fmt.Sprintf(imageFolderTemplate, folder, subFolder)
	)

	entries, err := assetsFS.ReadDir(folderAbsPath)
	if err != nil {
		return "", fmt.Errorf("read folder %q: %w", folderAbsPath, err)
	}

	switch len(entries) {
	case 0:
		return "", fmt.Errorf("empty folder: %q", folderAbsPath)

	case 1:
		return fmt.Sprintf("%s/%s", folderAbsPath, entries[0].Name()), nil
	}

	//nolint:gosec
	rndFileIdx := rand.Int() % len(entries)

	return fmt.Sprintf("%s/%s", folderAbsPath, entries[rndFileIdx].Name()), nil
}

// imageFolderAndSubfolder returns image folder and subfolder for shot type.
func imageFolderAndSubfolder(shot shotimage.ShotImage) (string, string) {
	switch shot.Type {
	case shotimage.ImageTypePrepareAttack:
		return prepareFolder, attackFolder

	case shotimage.ImageTypePrepareDefend:
		return prepareFolder, defendFolder

	case shotimage.ImageTypeAttack:
		return attackShotImage(attackFolder, shot)

	case shotimage.ImageTypeDefend:
		return attackShotImage(defendFolder, shot)

	case shotimage.ImageTypeWin:
		return resultFolder, winFolder

	case shotimage.ImageTypeLose:
		return resultFolder, loseFolder

	case shotimage.ImageTypeDraw:
		return resultFolder, drawFolder
	}

	return "", ""
}

// attackShotImage helper for imageFolderAndSubfolder for shotimage.ImageTypeAttack and shotimage.ImageTypeDefend.
func attackShotImage(folder string, shot shotimage.ShotImage) (string, string) {
	if shot.AttackerActualSide == game.ShotSideMiss {
		return folder,
			fmt.Sprintf("%s_%s_%s",
				shot.AttackerActualSide,
				shot.AttackerExpectedSide,
				shot.DefenderActualSide,
			)
	}

	return folder,
		fmt.Sprintf("%s_%s",
			shot.AttackerActualSide,
			shot.DefenderActualSide,
		)
}
