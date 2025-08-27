package order

import "github.com/korchizhinskiy/rocket-factory/order/internal/repository"

type service struct {
	repository repository.OrderRepository
}

func NewService(repository repository.OrderRepository) *service {
	return &service{
		repository: repository,
	}
}
