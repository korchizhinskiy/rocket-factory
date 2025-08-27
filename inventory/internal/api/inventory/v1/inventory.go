package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	responseConverter "github.com/korchizhinskiy/rocket-factory/inventory/internal/api/converter/request"
	"github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
	"github.com/korchizhinskiy/rocket-factory/inventory/internal/service/contract"
	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
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

func (api *api) ListPart(
	ctx context.Context,
	request *inventoryv1.ListPartRequest,
) (*inventoryv1.ListPartResponse, error) {

	var uuids []uuid.UUID
	if uuidFilter := request.GetFilter(); uuidFilter != nil {
		for _, id := range request.GetFilter().Uuids {
			parsedID, err := uuid.Parse(id)
			if err == nil {
				uuids = append(uuids, parsedID)
			}
		}
	}

	var categories []string
	if categoryFilter := request.GetFilter(); categoryFilter != nil {
		for _, category := range request.GetFilter().Categories {
			categories = append(categories, category.String())
		}
	}

	var tags []string
	if tagFilter := request.GetFilter(); tagFilter != nil {
		tags = append(tags, request.GetFilter().Tags...)
	}

	var names []string
	if nameFilter := request.GetFilter(); nameFilter != nil {
		names = append(names, request.GetFilter().Names...)
	}

	var manufacturerCountries []string
	if manufacturerCountriesFilter := request.GetFilter(); manufacturerCountriesFilter != nil {
		manufacturerCountries = append(manufacturerCountries, request.GetFilter().ManufactorerCountries...)
	}

	listFilter := contract.ListPartFilter{
		UUIDs:                 &uuids,
		Tags:                  &tags,
		Names:                 &names,
		Categories:            &categories,
		ManufactorerCountries: &manufacturerCountries,
	}

	partList, _ := api.inventoryService.GetListPart(ctx, listFilter)
	requestParts := make([]*inventoryv1.Part, 0, len(partList))
	for _, part := range partList {
		requestParts = append(requestParts, responseConverter.ConvertPartModelToGetPartResponse(part))
	}
	return &inventoryv1.ListPartResponse{Parts: requestParts}, nil
}
