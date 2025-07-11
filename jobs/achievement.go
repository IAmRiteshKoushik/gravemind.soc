package jobs

import (
	"context"
	"strconv"
	"time"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
	"github.com/redis/go-redis/v9"
)

func IncrDoc(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cmd.Valkey.HIncrBy(ctx, cmd.DocSet, username, 1).Err(); err != nil {
		return err
	}
	return nil
}

func CheckDoc(username string) (int, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	countStr, err := cmd.Valkey.HGet(ctx, cmd.DocSet, username).Result()
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

func IncrHelp(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cmd.Valkey.HIncrBy(ctx, cmd.HelpSet, username, 1).Err(); err != nil {
		return err
	}
	return nil
}

func CheckHelp(username string) (int, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	countStr, err := cmd.Valkey.HGet(ctx, cmd.HelpSet, username).Result()
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

func IncrTesting(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cmd.Valkey.HIncrBy(ctx, cmd.TestSet, username, 1).Err(); err != nil {
		return err
	}
	return nil
}

func CheckTesting(username string) (int, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	countStr, err := cmd.Valkey.HGet(ctx, cmd.TestSet, username).Result()
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

func IncrBugReport(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cmd.Valkey.HIncrBy(ctx, cmd.BugSet, username, 1).Err(); err != nil {
		return err
	}
	return nil
}

func CheckBugReport(username string) (int, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	countStr, err := cmd.Valkey.HGet(ctx, cmd.BugSet, username).Result()
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

func IncrFeature(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cmd.Valkey.HIncrBy(ctx, cmd.FeatSet, username, 1).Err(); err != nil {
		return err
	}
	return nil
}

func CheckFeature(username string) (int, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	countStr, err := cmd.Valkey.HGet(ctx, cmd.FeatSet, username).Result()
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

// -- Enamoured
// Get a PR accepted every week for a month
func CheckEnamoured(username string) {

}
