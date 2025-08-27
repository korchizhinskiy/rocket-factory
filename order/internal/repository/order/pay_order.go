package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
)

func (r *repository) PayOrder(ctx context.Context, orderID uuid.UUID, paymentMethod model.PaymentMethod) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, order := range r.orders {
		if order.ID == orderID {
			if *order.Status == model.OrderStatusPAID {
				return model.ErrOrderAlreadyPaid
			}
			order.Status = lo.ToPtr(model.OrderStatusPAID)
			order.PaymentMethod = lo.ToPtr(paymentMethod)
			return nil
		}
	}
	return model.ErrOrderNotFound
}
