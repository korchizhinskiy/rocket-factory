package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/payment/internal/service/dto"
)

func (s *service) Pay(ctx context.Context, payment dto.PaymentDTOIn) (uuid.UUID, error) {
	transactionID, err := s.paymentRepository.Pay(ctx, payment)
	if err != nil {
		return uuid.Nil, err
	}

	return transactionID, nil
}
