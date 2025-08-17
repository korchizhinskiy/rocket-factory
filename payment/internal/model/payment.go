package model

import "github.com/google/uuid"

type PaymentMethod byte

const (
	PAYMENT_METHOD_UNSPECIFIED = iota
	PAYMENT_METHOD_CARD
	PAYMENT_METHOD_SBP
	PAYMENT_METHOD_CREDIT_CARD
	PAYMENT_METHOD_INVESTOR_MONEY
)

type Payment struct {
	orderID       uuid.UUID
	userID        uuid.UUID
	paymentMethod PaymentMethod
}
type PaymentTransaction struct {
	ID uuid.UUID
	PaymentID    uuid.UUID
  PaymentMethod PaymentMethod
}
