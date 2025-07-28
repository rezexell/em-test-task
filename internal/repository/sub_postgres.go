package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubPostgres struct {
	pool *pgxpool.Pool
}

func NewSubPostgres(pool *pgxpool.Pool) *SubPostgres {
	return &SubPostgres{pool: pool}
}

func (r *SubPostgres) GetAllSubs() string {
	return "Success"
}
