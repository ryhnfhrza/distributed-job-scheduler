package controller

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/exception"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/helper"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/web"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/service"
)

type JobControllerImpl struct {
	JobService service.JobService
}

func NewJobController(jobService service.JobService) JobController {
	return &JobControllerImpl{
		JobService: jobService,
	}
}

func (controller *JobControllerImpl) Create(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	jobCreateRequest := web.JobCreateRequest{}
	helper.ReadFromRequestBody(request, &jobCreateRequest)

	jobResponse := controller.JobService.Create(request.Context(), jobCreateRequest)

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "Created",
		Data:   jobResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JobControllerImpl) Update(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	jobUpdateRequest := web.JobUpdateRequest{}
	helper.ReadFromRequestBody(request, &jobUpdateRequest)

	idString := param.ByName("id")
	id, err := strconv.Atoi(idString)
	helper.PanicIfError(err)

	jobUpdateRequest.Id = int64(id)

	jobResponse := controller.JobService.Update(request.Context(), jobUpdateRequest)

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   jobResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JobControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	idString := param.ByName("id")
	id, err := strconv.Atoi(idString)
	helper.PanicIfError(err)

	controller.JobService.Delete(request.Context(), int64(id))

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "job deleted",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JobControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	idString := param.ByName("id")
	id, err := strconv.Atoi(idString)
	helper.PanicIfError(err)

	jobResponse := controller.JobService.FindById(request.Context(), int64(id))

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   jobResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JobControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {

	filter, err := helper.ParseExpenseFilter(request.URL.Query())
	if err != nil {
		panic(exception.NewBadRequest(err.Error()))
	}

	jobResponse := controller.JobService.FindAll(request.Context(), filter)

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   jobResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JobControllerImpl) FindJobLog(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	idJobString := param.ByName("job_id")
	idJob, err := strconv.Atoi(idJobString)
	helper.PanicIfError(err)

	jobLogResponse := controller.JobService.FindJobLog(request.Context(), int64(idJob))

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   jobLogResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
