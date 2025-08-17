package v1

import (
	service "github.com/korchizhinskiy/rocket-factory/payment/internal/service"
	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
