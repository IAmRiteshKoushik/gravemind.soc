package workflows

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
	"github.com/redis/go-redis/v9"
)

type LiveUpdate struct {
	Username  string `json:"github_username"`
	Message   string `json:"message"`
	EventType string `json:"event_type"` // Bounty, Issue, Top-3
	Timestamp int64  `json:"time"`       // time is in unix.milliseconds for less size
}

func WriteToStream(username string, message string, eventType string) error {
	entry := LiveUpdate{
		Username:  username,
		Message:   message,
		EventType: eventType,
		Timestamp: time.Now().UnixMilli(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = cmd.Valkey.XAdd(ctx, &redis.XAddArgs{
		Stream: "live-update-stream",
		Values: map[string]interface{}{
			"data": string(data),
		},
	}).Err()
	if err != nil {
		return err
	}
	return nil
}
