package dto

import (
	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/payment/internal/model/enum"
)

type PaymentDTOIn struct {
	OrderID       uuid.UUID
	UserID        uuid.UUID
	PaymentMethod enum.PaymentMethod
}

