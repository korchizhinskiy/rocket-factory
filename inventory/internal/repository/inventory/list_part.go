package inventory

import (
	"context"

	model "github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
	"github.com/korchizhinskiy/rocket-factory/inventory/internal/repository/converter"
	"github.com/korchizhinskiy/rocket-factory/inventory/internal/service/contract"
)

func (r *repository) GetListPart(ctx context.Context, filters contract.ListPartFilter) ([]model.Part, error) {
	parts := make([]model.Part, 0, len(r.inventory))
	for _, part := range r.inventory {
		parts = append(parts, converter.ConvertPartRepoModelToModel(part))
	}
	return parts, nil
}
