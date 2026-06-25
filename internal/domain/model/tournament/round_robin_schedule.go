package tournament

type teamMarker struct {
	TeamIdx int
	IsGhost bool
}

func makeTeamMarkers(teamCount int) []teamMarker {
	//nolint:mnd
	hasGhost := teamCount%2 == 1
	if hasGhost {
		teamCount++
	}

	markers := make([]teamMarker, teamCount)

	if hasGhost {
		markers[teamCount-1] = teamMarker{
			IsGhost: true,
		}
	}

	for i := range markers {
		markers[i].TeamIdx = i
	}

	return markers
}

func isNeedToSwitchHomeAway(teamIdx int, roundNumber int) bool {
	if teamIdx == 0 && roundNumber%2 == 1 {
		return true
	}

	return false
}

func rotateTeamMarkers(markers []teamMarker) {
	lastTeamMarker := markers[len(markers)-1]
	for idx := len(markers) - 1; idx > 1; idx-- {
		markers[idx] = markers[idx-1]
	}
	markers[1] = lastTeamMarker
}

func makeRoundMatches(
	markers []teamMarker,
	roundNumber int,
) Round {
	var (
		teamCount     = len(markers)
		gamesPerRound = teamCount / 2
		roundMatches  = make(Round, 0, gamesPerRound)
	)

	for idx := range teamCount / 2 {
		var (
			homeMarker = markers[idx]
			awayMarker = markers[teamCount-idx-1]
			homeIdx    = homeMarker.TeamIdx
			awayIdx    = awayMarker.TeamIdx
		)

		if homeMarker.IsGhost || awayMarker.IsGhost {
			continue
		}

		if isNeedToSwitchHomeAway(homeIdx, roundNumber) {
			homeIdx, awayIdx = awayIdx, homeIdx
		}

		roundMatches = append(roundMatches, Match{
			Home: MatchTeamInfo{
				Idx: homeIdx,
			},
			Away: MatchTeamInfo{
				Idx: awayIdx,
			},
		})
	}

	return roundMatches
}

func RoundRobinSchedule(
	teamCount int,
) []Round {
	if teamCount == 0 {
		return nil
	}

	var (
		teamMarkers = makeTeamMarkers(teamCount)
		roundsCount = len(teamMarkers) - 1
		schedule    = make([]Round, 0, roundsCount)
	)

	for round := range roundsCount {
		roundMatches := makeRoundMatches(teamMarkers, round)

		schedule = append(schedule, roundMatches)

		rotateTeamMarkers(teamMarkers)
	}

	return schedule
}
