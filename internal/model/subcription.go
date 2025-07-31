package model

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a user's subscription
// @Description Subscription information
type Subscription struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ServiceName string     `gorm:"type:text;not null" json:"service_name" binding:"required,min=2,max=255"`
	MonthlyCost int        `gorm:"not null;check:monthly_cost>0" json:"monthly_cost" binding:"required,gt=0"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id" binding:"required,uuid4"`
	StartDate   time.Time  `gorm:"type:date;not null" json:"start_date" binding:"required,datetime=01/2006"`
	EndDate     *time.Time `gorm:"type:date" json:"end_date,omitempty" binding:"omitempty,datetime=01/2006"`
}

// SubReq represents subscription request data
// @Description Subscription request payload
type SubReq struct {
	ServiceName string `json:"service_name" binding:"required,min=2,max=255"`
	MonthlyCost int    `json:"monthly_cost" binding:"required,gt=0"`
	UserID      string `json:"user_id" binding:"required,uuid4"`
	StartDate   string `json:"start_date" binding:"required,datetime=01/2006"`
	EndDate     string `json:"end_date,omitempty" binding:"omitempty,datetime=01/2006"`
}
