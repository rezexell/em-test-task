package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/rezexell/em-test-task/internal/model"
	"github.com/rezexell/em-test-task/internal/repository"
	"time"
)

type Subscription interface {
	CreateSubscription(ctx context.Context, sub *model.SubReq) error
	GetSubscription(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *model.SubReq) error
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	ListAllSubscriptions(ctx context.Context) ([]*model.Subscription, error)
	ListSubscriptionsWithFilters(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error)
	TotalSubscriptionCost(ctx context.Context, userID *uuid.UUID, serviceName *string, periodStart, periodEnd time.Time) (int, error)
}
type Service struct {
	Subscription
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Subscription: NewSubService(repo.Subscription)}
}
