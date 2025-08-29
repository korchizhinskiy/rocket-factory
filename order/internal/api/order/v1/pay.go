package v1

import (
	"context"

	orderv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	request *orderv1.OrderPayRequest,
	params orderv1.PayOrderParams,
) (orderv1.PayOrderRes, error) {
	order, ok := o.storage.orders[params.OrderUUID.String()]
	if !ok {
		return &orderv1.NotFoundError{
			Code:    404,
			Message: "Order was not found.",
		}, nil
	}

	resp, err := o.payClient.PayOrder(
		ctx,
		&paymentv1.PayOrderRequest{
			OrderUuid: params.OrderUUID.String(),
			UserUuid:  req.UserUUID.String(),
			PaymentMethod: ConvertPaymentMethodStrToProto(
				req.PaymentMethod,
			),
		},
	)
	if err != nil {
		return &orderv1.InternalError{
			Code:    500,
			Message: "Internal Error",
		}, nil
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return &orderv1.InternalError{
			Code:    500,
			Message: "Internal Error",
		}, nil
	}

	(&order.Status).SetTo(orderv1.OrderStatusPAID)
	(&order.PatmentMethod).SetTo(req.PaymentMethod)

	return &orderv1.OrderPayResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
}
