package service

import (
	"context"

	"github.com/google/uuid"
)

type TransactionID = uuid.UUID

type InventoryService interface {
	Pay(ctx context.Context, payment model.Payment) (TransactionID, error)
}

