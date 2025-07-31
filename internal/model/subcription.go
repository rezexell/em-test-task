package model

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"time"

	"github.com/google/uuid"
)

// Subscription represents a user's subscription
// @Description Subscription information
type Subscription struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ServiceName  string     `gorm:"type:text;not null" json:"service_name" binding:"required,min=2,max=255"`
	MonthlyCost  int        `gorm:"not null;check:monthly_cost>0" json:"monthly_cost" binding:"required,gt=0"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id" binding:"required,uuid4"`
	StartDate    time.Time  `gorm:"type:date;not null" json:"-"`
	EndDate      *time.Time `gorm:"type:date" json:"-"`
	StartDateStr string     `gorm:"-" json:"start_date" binding:"required,datetime=01/2006"`
	EndDateStr   string     `gorm:"-" json:"end_date,omitempty" binding:"omitempty,datetime=01/2006"`
}

func (s *Subscription) AfterBind() error {
	startDate, err := time.Parse("01/2006", s.StartDateStr)
	if err != nil {
		return err
	}
	s.StartDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	if s.EndDateStr != "" {
		endDate, err := time.Parse("01/2006", s.EndDateStr)
		if err != nil {
			return err
		}
		firstOfNextMonth := time.Date(endDate.Year(), endDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		lastDay := firstOfNextMonth.AddDate(0, 0, -1)
		s.EndDate = &lastDay
	}
	return nil
}

func (s *Subscription) ToResponse() gin.H {
	response := gin.H{
		"id":           s.ID,
		"service_name": s.ServiceName,
		"monthly_cost": s.MonthlyCost,
		"user_id":      s.UserID,
		"start_date":   s.StartDate,
	}

	if s.EndDate != nil {
		response["end_date"] = s.EndDate
	}

	return response
}

func RegisterCustomBindings() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(SubscriptionStructLevelValidation, Subscription{})
	}
}

func SubscriptionStructLevelValidation(sl validator.StructLevel) {
	sub := sl.Current().Interface().(Subscription)

	if sub.EndDateStr != "" {
		start, err1 := time.Parse("01/2006", sub.StartDateStr)
		end, err2 := time.Parse("01/2006", sub.EndDateStr)

		if err1 == nil && err2 == nil && end.Before(start) {
			sl.ReportError(sub.EndDateStr, "end_date", "EndDate", "end_before_start", "")
		}
	}
}
