package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
)

type (
	OrderID       = uuid.UUID
	TransactionID = uuid.UUID
)

type OrderRepository interface {
	GetOrderByID(ctx context.Context, orderID OrderID) (model.Order, error)
	// CreateOrder(ctx context.Context, order model.Order) (OrderID, error)
	PayOrder(ctx context.Context, orderID OrderID) error
	// CancelOrder(ctx context.Context, orderID OrderID) error
}
