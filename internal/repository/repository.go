package repository

import (
	"context"
	"github.com/google/uuid"

	//"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rezexell/em-test-task/internal/model"
)

type Subscription interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error)
	ListAll(ctx context.Context) ([]*model.Subscription, error)
}

type Repository struct {
	Subscription
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{Subscription: NewSubPostgres(pool)}
}
