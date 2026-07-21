package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Raghunandan-79/pulsory/internal/config"
	"github.com/Raghunandan-79/pulsory/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WorkerService struct {
	db       *gorm.DB
	redis    *redis.Client
	regionID string
	workerID string
}

type WebsiteEvent struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type WebsiteCheckResult struct {
	WebsiteID      string
	URL            string
	ResponseTimeMs int
	Status         models.WebsiteStatus
	RegionID       string
}

func NewWorkerService(regionID, workerID string) *WorkerService {
	return &WorkerService{
		db:       config.DB,
		redis:    config.RedisClient,
		regionID: regionID,
		workerID: workerID,
	}
}

func (w *WorkerService) Start() error {
	// Create consumer group if it doesn't exist
	ctx := context.Background()
	streamName := "betteruptime:website"
	
	err := w.redis.XGroupCreateMkStream(ctx, streamName, w.regionID, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	// Continuous processing loop
	for {
		messages, err := w.readMessages()
		if err != nil {
			log.Printf("Error reading messages: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if len(messages) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		// Process messages concurrently
		results := w.processMessages(messages)

		// Save results to database
		if err := w.saveResults(results); err != nil {
			log.Printf("Error saving results: %v", err)
		}

		// Acknowledge messages
		messageIDs := make([]string, len(messages))
		for i, msg := range messages {
			messageIDs[i] = msg.ID
		}
		
		if err := config.XAckBulk(w.regionID, messageIDs); err != nil {
			log.Printf("Error acknowledging messages: %v", err)
		} else {
			log.Printf("Processed and acknowledged %d website checks", len(messages))
		}
	}
}

func (w *WorkerService) readMessages() ([]redis.XMessage, error) {
	ctx := context.Background()
	streamName := "betteruptime:website"

	results, err := w.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    w.regionID,
		Consumer: w.workerID,
		Streams:  []string{streamName, ">"},
		Count:    5,
		Block:    0,
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(results) == 0 || len(results[0].Messages) == 0 {
		return nil, nil
	}

	return results[0].Messages, nil
}

func (w *WorkerService) processMessages(messages []redis.XMessage) []WebsiteCheckResult {
	results := make([]WebsiteCheckResult, len(messages))
	
	for i, msg := range messages {
		// Extract website data from message
		var event WebsiteEvent
		if idVal, ok := msg.Values["id"].(string); ok {
			event.ID = idVal
		}
		if urlVal, ok := msg.Values["url"].(string); ok {
			event.URL = urlVal
		}

		// Check website status
		result := w.checkWebsite(event.URL, event.ID)
		results[i] = result
		
		log.Printf("Checked %s: %s (%dms)", event.URL, result.Status, result.ResponseTimeMs)
	}
	
	return results
}

func (w *WorkerService) checkWebsite(url, websiteID string) WebsiteCheckResult {
	startTime := time.Now()
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Make request
	resp, err := client.Get(url)
	duration := time.Since(startTime)
	
	result := WebsiteCheckResult{
		WebsiteID:      websiteID,
		URL:            url,
		ResponseTimeMs: int(duration.Milliseconds()),
		RegionID:       w.regionID,
	}
	
	if err != nil {
		result.Status = models.StatusDown
		log.Printf("%s is DOWN (%v)", url, err)
	} else {
		defer resp.Body.Close()
		
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			result.Status = models.StatusUp
			log.Printf("%s is UP (%dms)", url, duration.Milliseconds())
		} else {
			result.Status = models.StatusDown
			log.Printf("%s returned status %d", url, resp.StatusCode)
		}
	}
	
	return result
}

func (w *WorkerService) saveResults(results []WebsiteCheckResult) error {
	if len(results) == 0 {
		return nil
	}

	// Create website ticks in batch
	ticks := make([]models.WebsiteTick, len(results))
	for i, result := range results {
		ticks[i] = models.WebsiteTick{
			ResponseTimeMS: result.ResponseTimeMs,
			Status:         result.Status,
			RegionID:       uuid.MustParse(result.RegionID),
			WebsiteID:      uuid.MustParse(result.WebsiteID),
		}
	}

	// Use batch insert for better performance
	return w.db.CreateInBatches(ticks, 100).Error
}