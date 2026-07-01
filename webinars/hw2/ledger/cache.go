package ledger

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitCache() error {
	addr := getenvDefault("REDIS_ADDR", "localhost:6379")

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		client.Close()
		return err
	}

	redisClient = client
	log.Println("connected to Redis")
	return nil
}

func CloseCache() error {
	if redisClient == nil {
		return nil
	}

	return redisClient.Close()
}

func requireCache() (*redis.Client, error) {
	if redisClient == nil {
		return nil, errors.New("redis is not initialized")
	}

	return redisClient, nil
}

func invalidateReportSummaryCache() {
	if redisClient == nil {
		return
	}

	ctx := context.Background()

	err := redisClient.Del(ctx, reportSummaryCacheKey).Err()
	if err != nil {
		log.Printf("failed to invalidate report summary cache: %v", err)
	}
}