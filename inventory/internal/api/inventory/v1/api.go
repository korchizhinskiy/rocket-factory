package v1

import (
	service "github.com/korchizhinskiy/rocket-factory/payment/internal/service"
	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer

	// paymentService service.PaymentService
}

func NewAPI(inventoryService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}

