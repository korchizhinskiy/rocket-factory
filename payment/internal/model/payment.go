package model

import "github.com/google/uuid"

type PaymentMethod byte

const (
	PaymentMethodCard PaymentMethod = iota
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

func (pm *PaymentMethod) String() string {
	switch *pm {
	case PaymentMethodCard:
		return "Card"
	case PaymentMethodSBP:
		return "SBP"
	case PaymentMethodCreditCard:
		return "Credit Card"
	case PaymentMethodInvestorMoney:
		return "Investor Money"
	default:
		return "Unknown"
	}
}

type Payment struct {
	OrderID       uuid.UUID
	UserID        uuid.UUID
	PaymentMethod *PaymentMethod
	TransactionID *uuid.UUID
}
