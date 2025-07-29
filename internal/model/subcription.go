package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name" binding:"required,min=2,max=255"`
	MonthlyCost int        `json:"monthly_cost" binding:"required,gt=0"`
	UserID      uuid.UUID  `json:"user_id" binding:"required,uuid4"`
	StartDate   time.Time  `json:"start_date" binding:"datetime=01/2006"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubReq struct {
	ServiceName string `json:"service_name" binding:"required,min=2,max=255"`
	MonthlyCost int    `json:"monthly_cost" binding:"required,gt=0"`
	UserID      string `json:"user_id" binding:"required,uuid4"`
	StartDate   string `json:"start_date" binding:"required,datetime=01/2006"`
	EndDate     string `json:"end_date,omitempty" binding:"omitempty,datetime=01/2006"`
}
