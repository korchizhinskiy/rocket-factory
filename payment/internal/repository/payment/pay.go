package payment

import (
	"context"

	"github.com/google/uuid"

	repoModel "github.com/korchizhinskiy/rocket-factory/payment/internal/repository/model"
	"github.com/korchizhinskiy/rocket-factory/payment/internal/service/dto"
)

func (r *repository) Pay(ctx context.Context, payment dto.PaymentDTOIn) (uuid.UUID, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactionID := uuid.New()

	r.payments[transactionID.String()] = repoModel.Payment{
		OrderID:       payment.OrderID.String(),
		UserID:        payment.UserID.String(),
		PaymentMethod: payment.PaymentMethod.String(),
		TransactionID: transactionID.String(),
	}
	return transactionID, nil
}
