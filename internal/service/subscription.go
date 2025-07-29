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

func reqToModel(req *model.SubReq) (*model.Subscription, error) {
	var sub *model.Subscription
	var ed *time.Time

	id := uuid.New()
	usID, _ := uuid.Parse(req.UserID)
	sd, _ := time.Parse("01/2006", req.StartDate)
	if req.EndDate != "" {
		parsedEd, _ := time.Parse("01/2006", req.EndDate)
		ed = &parsedEd

		// Проверяем что end date после start date
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
