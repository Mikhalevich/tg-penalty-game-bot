package randomgenerator

import (
	"crypto/rand"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/playerprovider"
)

var (
	_ playerprovider.NameGenerator = (*RandomGenerator)(nil)
)

type RandomGenerator struct {
	prefix string
	length int
}

func New(
	prefix string,
	length int,
) *RandomGenerator {
	return &RandomGenerator{
		prefix: prefix,
		length: length,
	}
}

func (rg *RandomGenerator) GenerateName() string {
	return fmt.Sprintf("%s%s", rg.prefix, rand.Text()[:rg.length])
}
