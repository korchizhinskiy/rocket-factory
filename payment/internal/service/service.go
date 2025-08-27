package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/payment/internal/model"
)

type TransactionID = uuid.UUID

type PaymentService interface {
	Pay(ctx context.Context, payment model.Payment) (TransactionID, error)
}
