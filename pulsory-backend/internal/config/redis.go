package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {
	// Get Redis URL from environment
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	
	RedisClient = redis.NewClient(opt)
	
	// Test connection
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	
	log.Println("Redis connected successfully")
}

// Optional: Consumer group functions matching your TypeScript redisstream
func XReadGroup(consumerGroup, workerId string, count int64) ([]redis.XMessage, error) {
	ctx := context.Background()
	streamName := "betteruptime:website"
	
	results, err := RedisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: workerId,
		Streams:  []string{streamName, ">"},
		Count:    count,
		Block:    0,
	}).Result()
	
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, nil
	}
	
	return results[0].Messages, nil
}

func XAck(consumerGroup, eventId string) error {
	ctx := context.Background()
	streamName := "betteruptime:website"
	
	return RedisClient.XAck(ctx, streamName, consumerGroup, eventId).Err()
}

func XAckBulk(consumerGroup string, eventIds []string) error {
	ctx := context.Background()
	streamName := "betteruptime:website"
	
	pipe := RedisClient.Pipeline()
	for _, eventId := range eventIds {
		pipe.XAck(ctx, streamName, consumerGroup, eventId)
	}
	
	_, err := pipe.Exec(ctx)
	return err
}