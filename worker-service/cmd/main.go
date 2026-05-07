package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/app"
	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/config"
	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/repository"
	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/service"
)

func main() {
	envPath := filepath.Join("..", "internal", "config", ".env")

	if p := os.Getenv("CONFIG_PATH"); p != "" {
		envPath = p
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("warning: failed to load %s: %v", envPath, err)
	}

	db := app.NewDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("failed to close database:", err)
		}
	}()

	redisClient := app.NewRedis()
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Println("failed to close redis:", err)
		}
	}()

	jobRepository := repository.NewJobRepository(db)
	queueRepository := repository.NewQueueRepository(redisClient)

	smtpConfig := config.LoadSMTPConfig()

	workerService := service.NewWorkerService(
		jobRepository,
		queueRepository,
		smtpConfig,
	)

	ctx := context.Background()

	log.Println("worker service running...")

	for {
		err := workerService.Process(ctx)
		if err != nil {
			fmt.Println("worker error:", err)
		}
	}
}
