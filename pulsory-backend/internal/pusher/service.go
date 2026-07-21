package pusher

import (
	"context"
	"fmt"
	"log"

	"github.com/Raghunandan-79/pulsory/internal/config"
	"github.com/redis/go-redis/v9"
)

type Website struct {
	ID  string `gorm:"column:id"`
	URL string `gorm:"column:url"`
}

type PusherService struct {
	redis *redis.Client
}

func NewPusherService() *PusherService {
	return &PusherService{
		redis: config.RedisClient,
	}
}

func (p *PusherService) SyncWebsites() error {
	ctx := context.Background()
	
	// Fetch all websites with id and url only
	var websites []Website
	if err := config.DB.Table("websites").Select("id, url").Find(&websites).Error; err != nil {
		return fmt.Errorf("failed to fetch websites: %w", err)
	}
	
	if len(websites) == 0 {
		log.Println("No websites found to sync")
		return nil
	}
	
	// Stream name matching your TypeScript code
	streamName := "betteruptime:website"
	
	// Use pipeline for bulk operations
	pipe := p.redis.Pipeline()
	
	for _, website := range websites {
		// Add to Redis stream with the same structure as your TypeScript code
		pipe.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			ID:     "*", // Auto-generate ID like in your TypeScript
			Values: map[string]interface{}{
				"url": website.URL,
				"id":  website.ID,
			},
		})
	}
	
	// Execute all commands
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add to Redis stream: %w", err)
	}
	
	log.Printf("Synced %d websites to Redis stream '%s'", len(websites), streamName)
	return nil
}