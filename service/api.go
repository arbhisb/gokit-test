package service

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/junereycasuga/gokit-grpc-demo/repository"
)

type service struct {
	logger log.Logger
	db     repository.Repository
}

// Service interface
type Service interface {
	Add(ctx context.Context, numA, numB float32) (float32, error)
	Subtract(ctx context.Context, numA, numB float32) (float32, error)
	Multiply(ctx context.Context, numA, numB float32) (float32, error)
	Divide(ctx context.Context, numA, numB float32) (float32, error)
	Cda(ctx context.Context) (interface{}, error)
}

// NewService func initializes a service
func NewService(logger log.Logger, repo repository.Repository) Service {
	return &service{
		logger: logger,
		db:     repo,
	}
}

func (s service) Add(ctx context.Context, numA, numB float32) (float32, error) {
	return numA + numB, nil
}

func (s service) Subtract(ctx context.Context, numA, numB float32) (float32, error) {
	return numA - numB, nil
}

func (s service) Multiply(ctx context.Context, numA, numB float32) (float32, error) {
	return numA * numB, nil
}

func (s service) Divide(ctx context.Context, numA, numB float32) (float32, error) {
	return numA / numB, nil
}

func (s *service) Cda(ctx context.Context) (interface{}, error) {
	transactionEcv, err := s.db.GetTransaction2(ctx)
	if err != nil {
		return "", err
	}

	fmt.Println("transaction mybb : ", transactionEcv)
	return "success", nil
}
