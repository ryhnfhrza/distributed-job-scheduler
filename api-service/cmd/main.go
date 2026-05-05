package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/app"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/controller"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/exception"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/helper"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/repository"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/service"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/util"
)

func main() {
	envPath := filepath.Join("..", "internal", "config", ".env")

	if p := os.Getenv("CONFIG_PATH"); p != "" {
		envPath = p
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Failed to load %s: %v", envPath, err)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	DB := app.NewDB()

	validate := validator.New()
	util.NewValidator(validate)

	jobRepository := repository.NewJobRepository(DB)
	jobService := service.NewJobService(jobRepository, validate)
	jobController := controller.NewJobController(jobService)

	router := app.NewRouter(jobController)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    ":" + port,
		Handler: app.CORS(router),
	}

	log.Printf("Server running on port %s", port)
	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
