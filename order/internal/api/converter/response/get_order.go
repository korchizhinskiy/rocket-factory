package response

import (
	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
	orderv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/openapi/order/v1"
)

func ConvertModelToResponse(order model.Order) *orderv1.OrderDto {
	return &orderv1.OrderDto{
		OrderUUID:       order.ID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      orderv1.NewOptFloat64(*order.TotalPrice),
		TransactionUUID: orderv1.NewOptUUID(*order.TransactionUUID),
		PaymentMethod:   orderv1.NewOptPaymentMethod(orderv1.PaymentMethod(*order.PatmentMethod)),
		Status:          orderv1.NewOptOrderStatus(orderv1.OrderStatus(*order.Status)),
	}
}
