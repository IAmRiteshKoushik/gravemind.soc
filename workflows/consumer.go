package workflows

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
	"github.com/redis/go-redis/v9"
)

type BountyAction struct {
	ParticipantUsername string `json:"github_username"`
	Amount              int    `json:"amount"`
	Url                 string `json:"url"`
	Action              string `json:"action"`
}

type Achievement struct {
	ParticipantUsername string `json:"github_username"`
	Url                 string `json:"url"`
	Type                string `json:"type"`
}

type Solution struct {
	ParticipantUsername string `json:"github_username"`
	Url                 string `json:"pull_request_url"`
	Merged              bool   `json:"merged"`
}

func ReadBountyStream() {
	id, err := ReadLastEntry("bounty")
	if err != nil {
		cmd.Log.Error("Could not locate last entry", err)
		id = "$"
	}

	for {
		args := &redis.XReadArgs{
			Streams: []string{cmd.Bounty, id},
			Count:   1,
			Block:   0,
		}
		streams, err := cmd.Valkey.XRead(context.Background(), args).Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(10 * time.Second)
				continue
			}
			cmd.Log.Error("failed to read from bounty-stream. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Extract and process stream entries
		for _, stream := range streams {
			for _, message := range stream.Messages {

				// Mark event beginning in SQLite
				BountyEntry(message.ID)
				id = message.ID

				for _, val := range message.Values {

					var result BountyAction
					err := json.Unmarshal([]byte(val.(string)), &result)
					if err != nil {
						// Skip due to malformed JSON
						cmd.Log.Error("Failed to unmarshal JSON at bounty-stream", err)
						continue
					}
					BountyRunner(result)
				}
			}
		}
		// End of processing, reading next stream element
	}
}

func ReadAchivementStream() {
	id, err := ReadLastEntry("achievement")
	if err != nil {
		cmd.Log.Error("Could not locate last entry", err)
		id = "$"
	}

	for {
		args := &redis.XReadArgs{
			Streams: []string{cmd.SolutionMerge, id},
			Count:   1,
			Block:   0,
		}
		streams, err := cmd.Valkey.XRead(context.Background(), args).Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(10 * time.Second)
				continue
			}
			cmd.Log.Error("failed to read from achivement-stream. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Extract and process stream entries
		for _, stream := range streams {
			for _, message := range stream.Messages {

				// Mark event beginning in SQLite
				AchivementEntry(message.ID)
				id = message.ID

				for _, val := range message.Values {

					var result Achievement
					err := json.Unmarshal([]byte(val.(string)), &result)
					if err != nil {
						// Skip due to malformed JSON
						cmd.Log.Error("Failed to unmarshal JSON at achivement-stream", err)
						continue
					}
					AchievementRunner(result)
				}
			}
		}
		// End of processing, reading next stream element
	}
}

func ReadSolutionStream() {
	id, err := ReadLastEntry("pull_request")
	if err != nil {
		cmd.Log.Error("Could not locate last entry", err)
		id = "$"
	}

	for {
		args := &redis.XReadArgs{
			Streams: []string{cmd.SolutionMerge, id},
			Count:   1,
			Block:   0,
		}
		streams, err := cmd.Valkey.XRead(context.Background(), args).Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(10 * time.Second)
				continue
			}
			cmd.Log.Error("failed to read from solution-stream. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Extract and process stream entries
		for _, stream := range streams {
			for _, message := range stream.Messages {

				// Mark event beginning in SQLite
				PullRequestEntry(message.ID)
				id = message.ID

				for _, val := range message.Values {

					var result Solution
					err := json.Unmarshal([]byte(val.(string)), &result)
					if err != nil {
						// Skip due to malformed JSON
						cmd.Log.Error("Failed to unmarshal JSON at solution-stream", err)
						continue
					}
					PullRequestRunner(result)
				}
			}
		}
		// End of processing, reading next stream element
	}
}
