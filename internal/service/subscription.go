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

func (s *SubService) CreateSubscription(ctx context.Context, sub *model.Subscription) error {
	sub.ID = uuid.New()

	if sub.StartDate.IsZero() {
		sub.StartDate = time.Now().UTC()
	}

	if err := validateSubscription(sub); err != nil {
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

func (s *SubService) UpdateSubscription(ctx context.Context, sub *model.Subscription) error {
	if sub.ID == uuid.Nil {
		return errors.New("subscription ID is required for update")
	}

	if err := validateSubscription(sub); err != nil {
		return err
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

// validateSubscription валидирует данные подписки
func validateSubscription(sub *model.Subscription) error {
	if sub.ServiceName == "" {
		return errors.New("service name is required")
	}
	if sub.MonthlyCost <= 0 {
		return errors.New("monthly cost must be positive")
	}
	if sub.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if sub.StartDate.IsZero() {
		return errors.New("start date is required")
	}

	// Проверка, что дата окончания не раньше даты начала
	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return errors.New("end date cannot be before start date")
	}

	return nil
}
