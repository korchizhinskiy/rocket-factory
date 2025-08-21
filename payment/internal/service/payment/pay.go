package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, payment model.Payment) (uuid.UUID, error) {
	transactionID, err := s.paymentRepository.Pay(ctx, payment)
	if err != nil {
		return uuid.Nil, err
	}

	return transactionID, nil
}
