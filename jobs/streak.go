package jobs

import (
	"context"
	"strconv"
	"time"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
	"github.com/redis/go-redis/v9"
)

func AddStreak(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cmd.Valkey.HIncrBy(ctx, "stream", username, 1).Result()
	if err != nil {
		return err
	}
	// Setting up a 7 day TTL to handle the streak
	err = cmd.Valkey.Expire(ctx, "stream", 7*24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func CheckStreak(username string) (int, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	countStr, err := cmd.Valkey.HGet(ctx, cmd.EnamouredSet, username).Result()
	if err == redis.Nil {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, true, err
	}
	return count, true, nil
}
