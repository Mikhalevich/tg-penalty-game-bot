package model

import (
	"database/sql"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Player struct {
	ID                    int            `db:"id"`
	ChatID                int64          `db:"chat_id"`
	DisplayName           string         `db:"display_name"`
	CreatedAt             time.Time      `db:"created_at"`
	IsChangeNameTriggered bool           `db:"is_change_name_triggered"`
	NameChangedAt         sql.NullTime   `db:"name_changed_at"`
	GameStatus            string         `db:"game_status"`
	GameStatusChangedAt   sql.NullTime   `db:"game_status_changed_at"`
	CurrentGameID         sql.NullString `db:"current_game_id"`
	Score                 int            `db:"score"`
}

func (p Player) ToDomainPlayer() player.Player {
	return player.Player{
		ID:                    player.IDFromInt(p.ID),
		ChatID:                msginfo.ChatIDFromInt64(p.ChatID),
		DisplayName:           p.DisplayName,
		CreatedAt:             p.CreatedAt,
		IsChangeNameTriggered: p.IsChangeNameTriggered,
		NameChangedAt:         p.NameChangedAt.Time,
		GameStatus:            player.GameStatus(p.GameStatus),
		GameStatusChangedAt:   p.GameStatusChangedAt.Time,
		CurrentGameID:         p.CurrentGameID.String,
		Score:                 p.Score,
	}
}

func ToDomainPlayers(dbPlayers []Player) []player.Player {
	if len(dbPlayers) == 0 {
		return nil
	}

	players := make([]player.Player, 0, len(dbPlayers))

	for _, plr := range dbPlayers {
		players = append(players, plr.ToDomainPlayer())
	}

	return players
}

func ToDBPlayer(plr player.Player) Player {
	return Player{
		ID:                    plr.ID.Int(),
		ChatID:                plr.ChatID.Int64(),
		DisplayName:           plr.DisplayName,
		CreatedAt:             plr.CreatedAt,
		IsChangeNameTriggered: plr.IsChangeNameTriggered,
		GameStatus:            plr.GameStatus.String(),
		GameStatusChangedAt:   toNullTime(plr.GameStatusChangedAt),
		CurrentGameID:         toNullString(plr.CurrentGameID),
		Score:                 plr.Score,
	}
}

func ToDBPlayers(players []player.Player) []Player {
	dbPlayers := make([]Player, 0, len(players))

	for _, plr := range players {
		dbPlayers = append(dbPlayers, ToDBPlayer(plr))
	}

	return dbPlayers
}

func toNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func toNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
