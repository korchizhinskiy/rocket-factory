package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	responseConverter "github.com/korchizhinskiy/rocket-factory/inventory/internal/api/converter/request"
	"github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (api *api) GetPart(
	ctx context.Context,
	request *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	partUUID, _ := uuid.Parse(request.Uuid)

	part, err := api.inventoryService.GetPart(ctx, partUUID)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(
				codes.NotFound,
				"part with UUID %s not found", partUUID,
			)
		}
	}
	return &inventoryv1.GetPartResponse{Part: responseConverter.ConvertPartModelToGetPartResponse(part)}, nil
}
