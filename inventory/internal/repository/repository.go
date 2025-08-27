package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
	"github.com/korchizhinskiy/rocket-factory/inventory/internal/service/contract"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, partID uuid.UUID) (model.Part, error)
	GetListPart(ctx context.Context, filters contract.ListPartFilter) ([]model.Part, error)
}
