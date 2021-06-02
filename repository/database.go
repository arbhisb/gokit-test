package repository

import (
	"context"
)

type Transaction struct {
}

type Repository interface {
	GetTransaction(ctx context.Context, accountNo, postingDateStart, postingDateEnd string) ([]Transaction, error)
	GetTransaction2(ctx context.Context) ([]Transaction, error)
}
