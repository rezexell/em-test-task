package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/rezexell/em-test-task/internal/model"
	"github.com/rezexell/em-test-task/internal/repository"
)

type Subscription interface {
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
	GetSubscription(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *model.Subscription) error
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	ListUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error)
	ListAllSubscriptions(ctx context.Context) ([]*model.Subscription, error)
}
type Service struct {
	Subscription
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Subscription: NewSubService(repo.Subscription)}
}
