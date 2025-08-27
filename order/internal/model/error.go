package model

import "errors"

var (
	ErrOrderNotFound        = errors.New("order was not found in system")
	ErrOrderAlreadyPaid     = errors.New("order already paid")
	ErrOrderAlreadyCanceled = errors.New("order already canceled")
)
