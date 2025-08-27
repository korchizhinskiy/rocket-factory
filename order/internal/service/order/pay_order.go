package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
	"github.com/korchizhinskiy/rocket-factory/order/internal/service/payment"
)

func (s *service) PayOrder(
	ctx context.Context,
	paymentService payment.PaymentService,
	orderID uuid.UUID,
	paymentMethod model.PaymentMethod,
) (transactionID uuid.UUID, err error) {
	order, err := s.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.repository.PayOrder(ctx, orderID)
	if err != nil {
		return uuid.UUID{}, err
	}

	transactionID, err = paymentService.PayOrder(ctx, order)
	if err != nil {
		return uuid.UUID{}, err
	}

	return transactionID, nil
}
