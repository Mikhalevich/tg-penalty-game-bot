package lbrefresher_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/lbrefresher"
)

type LeaderboardRefresherSuite struct {
	*suite.Suite

	ctrl *gomock.Controller

	repo             *lbrefresher.MockRepository
	transactor       *lbrefresher.MockTransactor
	settingsProvider *lbrefresher.MockSettingsProvider
	timeProvider     *lbrefresher.MockTimeProvider

	refresher *lbrefresher.LeaderboardRefresher
}

func TestLeaderboardRefresherSuit(t *testing.T) {
	t.Parallel()
	suite.Run(t, &LeaderboardRefresherSuite{
		Suite: new(suite.Suite),
	})
}

func (s *LeaderboardRefresherSuite) SetupSuite() {
	s.ctrl = gomock.NewController(s.T())

	s.repo = lbrefresher.NewMockRepository(s.ctrl)
	s.transactor = lbrefresher.NewMockTransactor(s.ctrl)
	s.settingsProvider = lbrefresher.NewMockSettingsProvider(s.ctrl)
	s.timeProvider = lbrefresher.NewMockTimeProvider(s.ctrl)

	s.refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
}

func (s *LeaderboardRefresherSuite) TearDownSuite() {
}

func (s *LeaderboardRefresherSuite) TearDownTest() {
}

func (s *LeaderboardRefresherSuite) TearDownSubTest() {
}
