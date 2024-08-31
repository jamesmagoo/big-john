package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// store provide all functions to execute db queries and transactions, this extned the auto-gnerated Queries wit htxn capability
type Store struct {
	connPool *pgxpool.Pool
	*Queries 
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		connPool: connPool,
		Queries: New(connPool),
	}
}