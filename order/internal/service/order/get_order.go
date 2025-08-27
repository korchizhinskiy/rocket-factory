package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
)

func (s *service) GetOrderByID(ctx context.Context, orderID uuid.UUID) (model.Order, error) {
	return s.repository.GetOrderByID(ctx, orderID)
}
