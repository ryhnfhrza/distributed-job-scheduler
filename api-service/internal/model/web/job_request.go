package web

import "time"

type JobCreateRequest struct {
	Name    string      `json:"name" validate:"required,min=3,max=255"`
	Type    string      `json:"type" validate:"required,oneof=email webhook"`
	RunAt   *time.Time  `json:"run_at" validate:"required,future"`
	Payload interface{} `json:"payload" validate:"required"`
}

type JobUpdateRequest struct {
	Id    int64      `json:"id" validate:"required"`
	Name  string     `json:"name" validate:"omitempty,min=3,max=255"`
	RunAt *time.Time `json:"run_at" validate:"omitempty,future"`
}
