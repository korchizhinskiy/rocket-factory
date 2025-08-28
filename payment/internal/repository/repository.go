// Package repository provides an interface for payment operations.
package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/payment/internal/service/dto"
)

type PaymentRepository interface {
	Pay(ctx context.Context, payment dto.PaymentDTOIn) (uuid.UUID, error)
}
