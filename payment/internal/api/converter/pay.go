package converter

import (
	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/payment/internal/model/enum"
	"github.com/korchizhinskiy/rocket-factory/payment/internal/service/dto"
	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

func ConvertPayRequestToDTO(req *paymentv1.PayOrderRequest) dto.PaymentDTOIn {
	orderID, _ := uuid.Parse(req.OrderUuid) // nolint: gosec
	userID, _ := uuid.Parse(req.UserUuid)   // nolint: gosec

	return dto.PaymentDTOIn{
		OrderID:       orderID,
		UserID:        userID,
		PaymentMethod: enum.PaymentMethod(req.PaymentMethod),
	}
}
