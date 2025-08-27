package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
	"github.com/korchizhinskiy/rocket-factory/order/internal/service/payment"
)

type transactionID = uuid.UUID

type OrderService interface {
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (model.Order, error)
	PayOrder(ctx context.Context, paymentService payment.PaymentService, orderID uuid.UUID) (transactionID, error)
}
