package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rezexell/em-test-task/internal/model"
)

type SubPostgres struct {
	pool *pgxpool.Pool
}

func NewSubPostgres(pool *pgxpool.Pool) *SubPostgres {
	return &SubPostgres{pool: pool}
}

func (r *SubPostgres) Create(ctx context.Context, sub *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, service_name, monthly_cost, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(ctx, query,
		sub.ID,
		sub.ServiceName,
		sub.MonthlyCost,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)

	return err
}

// GetByID возвращает подписку по ID
func (r *SubPostgres) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	query := `
		SELECT id, service_name, monthly_cost, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`

	var sub model.Subscription
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.MonthlyCost,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

// Update обновляет существующую подписку
func (r *SubPostgres) Update(ctx context.Context, sub *model.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $2, monthly_cost = $3, user_id = $4, start_date = $5, end_date = $6
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		sub.ID,
		sub.ServiceName,
		sub.MonthlyCost,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

// Delete удаляет подписку по ID
func (r *SubPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

// ListByUser возвращает все подписки для указанного пользователя
func (r *SubPostgres) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error) {
	query := `
		SELECT id, service_name, monthly_cost, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY start_date DESC
	`

	return r.querySubscriptions(ctx, query, userID)
}

// ListAll возвращает все подписки
func (r *SubPostgres) ListAll(ctx context.Context) ([]*model.Subscription, error) {
	query := `
		SELECT id, service_name, monthly_cost, user_id, start_date, end_date
		FROM subscriptions
		ORDER BY start_date DESC
	`

	return r.querySubscriptions(ctx, query)
}

func (r *SubPostgres) querySubscriptions(ctx context.Context, query string, args ...interface{}) ([]*model.Subscription, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.MonthlyCost,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}
