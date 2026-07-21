package main

import (
	"log"
	"time"

	"github.com/Raghunandan-79/pulsory/internal/config"
	"github.com/Raghunandan-79/pulsory/internal/pusher"
)

func main() {
	// Connect to database
	config.ConnectDB()
	
	// Connect to Redis
	config.ConnectRedis()
	
	// Create pusher service
	pusherService := pusher.NewPusherService()
	
	// Run initial sync
	if err := pusherService.SyncWebsites(); err != nil {
		log.Printf("Initial sync error: %v", err)
	}
	
	// Schedule sync every 3 minutes
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	
	go func() {
		for range ticker.C {
			if err := pusherService.SyncWebsites(); err != nil {
				log.Printf("Sync error: %v", err)
			}
		}
	}()
	
	log.Println("Pusher started successfully")
	
	// Keep the program running
	select {}
}