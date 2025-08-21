// Package model defines the data structures used in the payment repository.
package model

type PaymentMethod byte

type Payment struct {
	OrderID       string
	UserID        string
	PaymentMethod string
	TransactionID string
}
