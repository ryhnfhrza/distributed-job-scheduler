package exception

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/helper"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/web"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {

	if notFoundError(writer, request, err) {

		return
	}

	if validationError(writer, request, err) {
		return
	}
	if badRequestError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}

func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exeception, ok := err.(*NotFoundError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   exeception.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {

		return false
	}
}

func validationError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exeception, ok := err.(validator.ValidationErrors)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   exeception.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {

		return false
	}
}

func badRequestError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(*BadRequest)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data: map[string]string{
				"message": exception.Message,
			},
		}
		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {

		return false
	}
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	webResponse := web.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data: map[string]string{
			"message": fmt.Sprintf("%v", err),
		},
	}

	helper.WriteToResponseBody(writer, webResponse)
}
