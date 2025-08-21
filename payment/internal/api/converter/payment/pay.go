package payment

import (
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/korchizhinskiy/rocket-factory/payment/internal/model"
	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

func ConvertPayRequestToModel(req *paymentv1.PayOrderRequest) model.Payment {
	orderID, _ := uuid.Parse(req.OrderUuid) // nolint: gosec
	userID, _ := uuid.Parse(req.UserUuid)   // nolint: gosec

	return model.Payment{
		OrderID:       orderID,
		UserID:        userID,
		PaymentMethod: lo.ToPtr(model.PaymentMethod(req.PaymentMethod)),
	}
}
