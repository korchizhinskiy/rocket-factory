package inventory

import (
	"sync"

	repoModel "github.com/korchizhinskiy/rocket-factory/inventory/internal/repository/model"
)

type repository struct {
	mu        sync.RWMutex
	inventory map[string]repoModel.Part
}

func NewRepository() *repository {
	inv := GenerateParts(1000)
	return &repository{
		inventory: inv,
	}
}
