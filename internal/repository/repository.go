package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Subscription interface {
	GetAllSubs() string

	//CreateSub()
	//GetSub()
	//UpdateSub()
	//DeleteSub()
}

type Repository struct {
	Subscription
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{Subscription: NewSubPostgres(pool)}
}
