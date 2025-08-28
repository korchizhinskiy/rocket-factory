package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/payment/internal/service/dto"
)

type TransactionID = uuid.UUID

type PaymentService interface {
	Pay(ctx context.Context, payment dto.PaymentDTOIn) (TransactionID, error)
}
