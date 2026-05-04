package helper

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/domain"
)

func ParseExpenseFilter(query url.Values) (domain.Filter, error) {
	var filter domain.Filter

	if limitStr := query.Get("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 0 {
			return filter, fmt.Errorf("invalid limit")
		}
		filter.Limit = l
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil || o < 0 {
			return filter, fmt.Errorf("invalid offset")
		}
		filter.Offset = o
	}
	if Status := query.Get("description"); Status != "" {
		filter.Status = Status
	}

	return filter, nil
}
