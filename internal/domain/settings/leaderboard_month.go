package settings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const (
	leaderboardMonth = "leaderboard_month"
)

type leaderboardMonthPayload struct {
	Month int `json:"month"`
}

func (s *Settings) GetLeaderboardMonth(ctx context.Context) (int, error) {
	item, err := s.repo.GetSetting(ctx, leaderboardMonth)
	if err != nil {
		return 0, fmt.Errorf("get setting: %w", err)
	}

	var payload leaderboardMonthPayload
	if err := json.NewDecoder(bytes.NewReader(item.Payload)).Decode(&payload); err != nil {
		return 0, fmt.Errorf("decode setting payload: %w", err)
	}

	return payload.Month, nil
}

func (s *Settings) SetLeaderboardMonth(ctx context.Context, month int) error {
	payload, err := json.Marshal(leaderboardMonthPayload{
		Month: month,
	})

	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	now := s.timeProvider.Now()

	if err := s.repo.SetSetting(
		ctx,
		SettingItem{
			ID:               leaderboardMonth,
			IsEnabled:        true,
			CreatedAt:        now,
			Payload:          payload,
			PayloadUpdatedAt: now,
		},
	); err != nil {
		return fmt.Errorf("set setting: %w", err)
	}

	return nil
}
