package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/ryhnfhrza/distributed-job-scheduler/scheduler-service/internal/app"
	"github.com/ryhnfhrza/distributed-job-scheduler/scheduler-service/internal/repository"
	"github.com/ryhnfhrza/distributed-job-scheduler/scheduler-service/internal/service"
)

func main() {
	envPath := filepath.Join("..", "internal", "config", ".env")

	if p := os.Getenv("CONFIG_PATH"); p != "" {
		envPath = p
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Failed to load %s: %v", envPath, err)
	}

	//DB
	DB := app.NewDB()

	// Redis
	redisClient := app.NewRedis()
	defer func() {
		if err := redisClient.Close(); err != nil {
			fmt.Println("Error closing Redis connection:", err)
		}
	}()

	// repository
	jobRepository := repository.NewJobRepository(DB)
	queueRepository := repository.NewQueueRepository(redisClient)

	//service
	schedulerService := service.NewSchedulerService(jobRepository, queueRepository)

	// Scheduler Loop
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	ctx := context.Background()

	log.Println("Scheduler started...")

	for range ticker.C {
		err := schedulerService.Process(ctx)
		if err != nil {
			log.Println("Scheduler error:", err)
		}
	}
}
