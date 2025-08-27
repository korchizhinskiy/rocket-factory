package inventory

import (
	"context"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, partID uuid.UUID) (model.Part, error) {
	return s.inventoryRepository.GetPart(ctx, partID)
}
