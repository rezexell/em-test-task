package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rezexell/em-test-task/internal/model"
	"github.com/rezexell/em-test-task/internal/repository"
	"time"
)

type SubService struct {
	repo repository.Subscription
}

func NewSubService(repo repository.Subscription) *SubService {
	return &SubService{repo: repo}
}

func (s *SubService) CreateSubscription(ctx context.Context, req *model.SubReq) error {
	sub, err := reqToModel(req)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, sub)
}

func (s *SubService) GetSubscription(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid subscription ID")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *SubService) UpdateSubscription(ctx context.Context, req *model.SubReq) error {
	sub, err := reqToModel(req)
	if err != nil {
		return err
	}
	if sub.ID == uuid.Nil {
		return errors.New("subscription ID is required for update")
	}

	return s.repo.Update(ctx, sub)
}

func (s *SubService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid subscription ID")
	}

	return s.repo.Delete(ctx, id)
}

func (s *SubService) ListUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.ListByUser(ctx, userID)
}

func (s *SubService) ListAllSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	return s.repo.ListAll(ctx)
}

func (s *SubService) TotalSubscriptionCost(ctx context.Context, userID *uuid.UUID, serviceName *string, periodStart, periodEnd time.Time) (int, error) {
	if periodStart.After(periodEnd) {
		return 0, errors.New("start period cannot be after end period")
	}

	subscriptions, err := s.repo.ListWithFilters(
		ctx,
		userID,
		serviceName,
		periodStart,
		periodEnd,
	)

	if err != nil {
		return 0, err
	}

	total := 0
	for _, sub := range subscriptions {
		activeMonths := calculateActiveMonths(
			sub.StartDate,
			sub.EndDate,
			periodStart,
			periodEnd,
		)

		total += sub.MonthlyCost * activeMonths
	}

	return total, nil
}

func calculateActiveMonths(subStart time.Time, subEnd *time.Time, periodStart, periodEnd time.Time) int {
	activityStart := subStart
	if subStart.Before(periodStart) {
		activityStart = periodStart
	}

	var activityEnd time.Time
	if subEnd == nil {
		activityEnd = periodEnd
	} else {
		activityEnd = *subEnd
		if activityEnd.After(periodEnd) {
			activityEnd = periodEnd
		}
	}

	if activityStart.After(activityEnd) {
		return 0
	}

	return countFullMonthsBetween(activityStart, activityEnd)
}

func countFullMonthsBetween(start, end time.Time) int {
	start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.UTC)

	months := 0
	for current := start; !current.After(end); current = current.AddDate(0, 1, 0) {
		months++
	}

	return months
}

func reqToModel(req *model.SubReq) (*model.Subscription, error) {
	var sub *model.Subscription
	var ed *time.Time

	id := uuid.New()
	usID, _ := uuid.Parse(req.UserID)
	sd, _ := time.Parse("01/2006", req.StartDate)
	if req.EndDate != "" {
		parsedEd, _ := time.Parse("01/2006", req.EndDate)
		ed = &parsedEd

		if ed.Before(sd) {
			return nil, errors.New("end date cannot be before start date")
		}
	}

	sub = &model.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		MonthlyCost: req.MonthlyCost,
		UserID:      usID,
		StartDate:   sd,
		EndDate:     ed,
	}

	return sub, nil
}
