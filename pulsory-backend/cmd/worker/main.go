package main

import (
	"log"
	"os"

	"github.com/Raghunandan-79/pulsory/internal/config"
	"github.com/Raghunandan-79/pulsory/internal/worker"
)

func main() {
	// Get region and worker IDs from environment
	regionID := os.Getenv("REGION_ID")
	workerID := os.Getenv("WORKER_ID")

	if regionID == "" {
		log.Fatal("REGION_ID environment variable is required")
	}
	if workerID == "" {
		log.Fatal("WORKER_ID environment variable is required")
	}

	// Connect to database
	config.ConnectDB()

	// Connect to Redis
	config.ConnectRedis()

	// Create worker service
	workerService := worker.NewWorkerService(regionID, workerID)

	// Start the worker
	log.Printf("🚀 Worker started: Region=%s, Worker=%s", regionID, workerID)
	
	if err := workerService.Start(); err != nil {
		log.Fatalf("Worker error: %v", err)
	}
}