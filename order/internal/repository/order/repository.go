package order

import (
	"sync"

	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
)

type repository struct {
	mu     sync.RWMutex
	orders []model.Order
}

func NewRepository() *repository {
  return &repository{
    orders: make([]model.Order, 0),
  }
}
