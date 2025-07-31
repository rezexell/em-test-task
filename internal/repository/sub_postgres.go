package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rezexell/em-test-task/internal/model"
	"gorm.io/gorm"
)

type SubPostgres struct {
	db *gorm.DB
}

func NewSubPostgres(db *gorm.DB) *SubPostgres {
	return &SubPostgres{db: db}
}

func (r *SubPostgres) Create(ctx context.Context, sub *model.Subscription) error {
	result := r.db.WithContext(ctx).Create(sub)
	return result.Error
}

func (r *SubPostgres) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	var sub model.Subscription
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&sub)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &sub, nil
}

func (r *SubPostgres) Update(ctx context.Context, sub *model.Subscription) error {
	result := r.db.WithContext(ctx).Model(&model.Subscription{}).
		Where("id = ?", sub.ID).
		Updates(sub)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

func (r *SubPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Subscription{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (r *SubPostgres) ListAll(ctx context.Context) ([]*model.Subscription, error) {
	var subscriptions []*model.Subscription
	result := r.db.WithContext(ctx).
		Order("start_date DESC").
		Find(&subscriptions)

	if result.Error != nil {
		return nil, result.Error
	}
	return subscriptions, nil
}

func (r *SubPostgres) ListWithFilters(ctx context.Context, userID *uuid.UUID, serviceName *string, startPeriod, endPeriod *time.Time) ([]*model.Subscription, error) {
	var subscriptions []*model.Subscription

	query := r.db.WithContext(ctx)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if serviceName != nil {
		query = query.Where("service_name = ?", *serviceName)
	}

	if startPeriod != nil && endPeriod != nil {
		query = query.Where("start_date <= ?", endPeriod).
			Where("(end_date IS NOT NULL AND end_date >= ?) OR (end_date IS NULL)", startPeriod)
	}

	result := query.Order("start_date DESC").Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
