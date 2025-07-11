package jobs

import (
	"context"
	"time"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
	"github.com/redis/go-redis/v9"
)

// Function which can check whether the Pirates of Issuebian badge is to be
// given out or not
func CheckIssuebian(username string) (bool, error) {
	sets := []string{
		cmd.CppRank,
		cmd.JavaRank,
		cmd.PyRank,
		cmd.JsRank,
		cmd.GoRank,
		cmd.RustRank,
		cmd.ZigRank,
		cmd.FlutterRank,
		cmd.KotlinRank,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, set := range sets {
		prMergeCount, err := cmd.Valkey.ZScore(ctx, set, username).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return false, err
		}
		if prMergeCount >= 10.0 {
			return true, nil
		}
	}
	return false, nil
}

func FindPrCount(username string, board string) (int, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	prMergeCount, err := cmd.Valkey.ZScore(ctx, board, username).Result()
	if err == redis.Nil {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return int(prMergeCount), true, nil
}

func IncrPrCount(username string, board string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First increment in language board
	if err := cmd.Valkey.ZIncrBy(ctx, board, 1, username).Err(); err != nil {
		return err
	}
	// Then increment in leaderboard
	if err := cmd.Valkey.ZIncrBy(ctx, cmd.Leaderboard, 0.001, username).Err(); err != nil {
		return err
	}
	return nil
}

func FindRank(username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rank, err := cmd.Valkey.ZRank(ctx, cmd.Leaderboard, username).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return int(rank), nil
}
