package redisbr

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/redis/go-redis/v9"
)

func (r *RedisButtonRepository) GetButton(ctx context.Context, id button.ID) (*button.Button, error) {
	key, btnNum := parseButtonID(id)
	if btnNum == "" {
		btn, err := r.singleButton(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("single button: %w", err)
		}

		return btn, nil
	}

	btn, err := r.hmapButton(ctx, key, btnNum)
	if err != nil {
		return nil, fmt.Errorf("hmap button: %w", err)
	}

	return btn, nil
}

func (r *RedisButtonRepository) singleButton(ctx context.Context, key string) (*button.Button, error) {
	blob, err := r.client.GetDel(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, perror.NotFound("single button not found")
		}

		return nil, fmt.Errorf("redis get: %w", err)
	}

	btn, err := decodeButton(blob)
	if err != nil {
		return nil, fmt.Errorf("decode button: %w", err)
	}

	return btn, nil
}

func (r *RedisButtonRepository) hmapButton(ctx context.Context, key, field string) (*button.Button, error) {
	blob, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, perror.NotFound("hmap not found")
		}

		return nil, fmt.Errorf("hget: %w", err)
	}

	btn, err := decodeButton([]byte(blob))
	if err != nil {
		return nil, fmt.Errorf("decode button: %w", err)
	}

	return btn, nil
}
