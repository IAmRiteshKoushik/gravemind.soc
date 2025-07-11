package jobs

import (
	"context"
	"time"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
)

func UpdateBounty(username string, amt int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := cmd.Valkey.ZIncrBy(ctx, cmd.Leaderboard, float64(amt), username).Err()
	if err != nil {
		return err
	}
	return nil
}
