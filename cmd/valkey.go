package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var Valkey *redis.Client

func InitValkey() (*redis.Client, error) {
	host := App.RedisHost
	port := App.RedisPort
	resp := 3

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0, // default DB
		Protocol: resp,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result() // health-check
	if err != nil {
		Log.Error("[FAIL]: Health-check failed for Valkey.", err)
		return nil, err
	}
	Log.Info(
		fmt.Sprintf("[PASSED]: Health-check successfuly for Valkey. Response: %s", pong))

	return rdb, nil
}

func CloseValkey(client *redis.Client) {
	if client != nil {
		client.Close()
	}
}
