package payment

import (
	"sync"

	repoModel "github.com/korchizhinskiy/rocket-factory/payment/internal/repository/model"
)

type repository struct {
	mu       sync.RWMutex
	payments map[string]repoModel.Payment
}

func NewRepository() *repository {
	return &repository{
		payments: make(map[string]repoModel.Payment),
	}
}
