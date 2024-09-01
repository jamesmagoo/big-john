package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	// TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	// CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	// VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// store provides all functions to execute db queries and transactions, this extends the auto-generated sqlc Queries with txn capability
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries 
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries: New(connPool),
	}
}