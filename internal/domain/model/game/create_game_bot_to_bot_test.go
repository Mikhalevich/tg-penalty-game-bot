//nolint:testpackage
package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdjustBotParams(t *testing.T) {
	t.Parallel()
	t.Run("easy to easy", func(t *testing.T) {
		t.Parallel()

		var (
			homeBot = createBot(BotDifficultyEasy)
			awayBot = createBot(BotDifficultyEasy)
		)

		adjustBotParams(&homeBot, &awayBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botEasy,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent20,
				DefendMissPercent: percent10,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}, homeBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botEasy,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent20,
				DefendMissPercent: percent10,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}, awayBot)
	})

	t.Run("easy to insane", func(t *testing.T) {
		t.Parallel()

		var (
			homeBot = createBot(BotDifficultyEasy)
			awayBot = createBot(BotDifficultyInsane)
		)

		adjustBotParams(&homeBot, &awayBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botEasy,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent20,
				DefendMissPercent: percent10,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}, homeBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botInsane,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent2,
				DefendMissPercent: percent2,
				ForceGoalPercent:  percent60,
				ForceSavePercent:  percent60,
			},
		}, awayBot)
	})

	t.Run("insane to easy", func(t *testing.T) {
		t.Parallel()

		var (
			homeBot = createBot(BotDifficultyInsane)
			awayBot = createBot(BotDifficultyEasy)
		)

		adjustBotParams(&homeBot, &awayBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botInsane,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent2,
				DefendMissPercent: percent2,
				ForceGoalPercent:  percent60,
				ForceSavePercent:  percent60,
			},
		}, homeBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botEasy,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent20,
				DefendMissPercent: percent10,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}, awayBot)
	})

	t.Run("hard to insane", func(t *testing.T) {
		t.Parallel()

		var (
			homeBot = createBot(BotDifficultyHard)
			awayBot = createBot(BotDifficultyInsane)
		)

		adjustBotParams(&homeBot, &awayBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botHard,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent2,
				DefendMissPercent: percent2,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}, homeBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botInsane,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent2,
				DefendMissPercent: percent2,
				ForceGoalPercent:  30,
				ForceSavePercent:  35,
			},
		}, awayBot)
	})

	t.Run("insave to insane", func(t *testing.T) {
		t.Parallel()

		var (
			homeBot = createBot(BotDifficultyInsane)
			awayBot = createBot(BotDifficultyInsane)
		)

		adjustBotParams(&homeBot, &awayBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botInsane,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent2,
				DefendMissPercent: percent2,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}, homeBot)

		require.Equal(t, Player{
			ID:          0,
			DisplayName: botInsane,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent2,
				DefendMissPercent: percent2,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}, awayBot)
	})
}
