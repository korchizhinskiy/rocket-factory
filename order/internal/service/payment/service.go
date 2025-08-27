package payment

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

type PaymentService struct {
	paymentClient paymentv1.PaymentServiceClient
}

func NewPaymentService(paymentClient paymentv1.PaymentServiceClient) *PaymentService {
	return &PaymentService{
		paymentClient: paymentClient,
	}
}

func (ps *PaymentService) PayOrder(ctx context.Context, order model.Order) (uuid.UUID, error) {
	resp, err := ps.paymentClient.PayOrder(
		ctx,
		&paymentv1.PayOrderRequest{
			OrderUuid:     order.ID.String(),
			UserUuid:      order.UserUUID.String(),
			PaymentMethod: ConvertPaymentMethodStrToProto(*order.PaymentMethod),
		},
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return uuid.UUID{}, err
	}

	return transactionUUID, nil
}
