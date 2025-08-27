package inventory

import (
	"context"

	"github.com/google/uuid"

	model "github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
	"github.com/korchizhinskiy/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(ctx context.Context, partID uuid.UUID) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.inventory[partID.String()]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}
	return converter.ConvertPartRepoModelToModel(part), nil
}
