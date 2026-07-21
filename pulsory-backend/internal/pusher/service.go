package pusher

import (
	"context"
	"encoding/json"
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
	
	// Prepare data for Redis stream using pipeline
	streamKey := "websites:stream"
	pipe := p.redis.Pipeline()
	
	for _, website := range websites {
		// Create the data payload matching your TypeScript structure
		data := map[string]interface{}{
			"id":  website.ID,
			"url": website.URL,
		}
		
		// Convert to JSON string (or use map directly)
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to marshal website %s: %v", website.ID, err)
			continue
		}
		
		// Add to Redis stream
		pipe.XAdd(ctx, &redis.XAddArgs{
			Stream: streamKey,
			Values: map[string]interface{}{
				"data": string(jsonData),
			},
		})
	}
	
	// Execute all commands in pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add to Redis stream: %w", err)
	}
	
	log.Printf("Synced %d websites to Redis stream", len(websites))
	return nil
}