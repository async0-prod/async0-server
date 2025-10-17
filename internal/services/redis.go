package services

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	opt, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
	client := redis.NewClient(opt)

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("Connected to Redis!")
	return client, nil
}
