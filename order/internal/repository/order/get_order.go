package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
)

func (r *repository) GetOrderByID(ctx context.Context, orderID uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, order := range r.orders {
		if order.ID == orderID {
			return order, nil
		}
	}

	return model.Order{}, model.ErrOrderNotFound
}
