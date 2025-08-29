package v1

import (
	"github.com/korchizhinskiy/rocket-factory/order/internal/service"
	orderv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

var _ orderv1.Handler = (*api)(nil)

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
