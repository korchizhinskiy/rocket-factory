package v1

import (
	"context"

	"github.com/korchizhinskiy/rocket-factory/payment/internal/api/converter/payment"
	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

func (api *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	transactionID, err := api.paymentService.Pay(ctx, payment.ConvertPayRequestToModel(req))
	if err != nil {
		return nil, err
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionID.String(),
	}, nil
}
