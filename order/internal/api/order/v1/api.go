package v1

import (
	"github.com/korchizhinskiy/rocket-factory/order/internal/service"
	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

type api struct {
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient

	orderService service.OrderService
}

func NewAPI(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	orderService service.OrderService,
) *api {
	return &api{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderService:    orderService,
	}
}
