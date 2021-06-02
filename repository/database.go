package repository

import (
	"context"
)

type Transaction struct {
	Id       string
	Offer_id string
	Quote_id string
}

type TransactionJson *struct {
	Id       string `json:"id"`
	Offer_id string `json:"offer_id"`
	Quote_id string `json:"quote_id"`
}

type Repository interface {
	GetTransaction(ctx context.Context, accountNo, postingDateStart, postingDateEnd string) ([]Transaction, error)
	GetTransaction2(ctx context.Context) ([]Transaction, error)
}
