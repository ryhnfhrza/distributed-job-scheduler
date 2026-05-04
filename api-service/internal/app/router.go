package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/controller"
)

func NewRouter(jobController controller.JobController) *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/v1/jobs", jobController.Create)
	router.PUT("/api/v1/jobs/:id", jobController.Update)
	router.DELETE("/api/v1/jobs/:id", jobController.Delete)
	router.GET("/api/v1/jobs/:id", jobController.FindById)
	router.GET("/api/v1/jobs", jobController.FindAll)
	router.GET("/api/v1/job_logs/:job_id", jobController.FindJobLog)

	return router
}
