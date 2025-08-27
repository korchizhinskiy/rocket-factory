package converter

import (
	"github.com/google/uuid"

	model "github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
	repoModel "github.com/korchizhinskiy/rocket-factory/inventory/internal/repository/model"
)

func ConvertPartRepoModelToModel(partRepoModel repoModel.Part) model.Part {
	partID, _ := uuid.Parse(partRepoModel.ID)
	modelMetadata := make(map[string]model.MetadataValue)
	for k, v := range *partRepoModel.Metadata {
		modelMetadata[k] = model.MetadataValue(
			model.MetadataValue{
				FloatValue:  v.FloatValue,
				IntValue:    v.IntValue,
				StringValue: v.StringValue,
				BoolValue:   v.BoolValue,
			},
		)
	}
	return model.Part{
		ID:          partID,
		Name:        partRepoModel.Name,
		Description: partRepoModel.Description,
		Price:       partRepoModel.Price,
		Category:    partRepoModel.Category,
		Dimensions: model.Dimension{
			Length: partRepoModel.Dimensions.Length,
			Width:  partRepoModel.Dimensions.Width,
			Height: partRepoModel.Dimensions.Height,
			Weight: partRepoModel.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    partRepoModel.Manufacturer.Name,
			Country: partRepoModel.Manufacturer.Country,
			Website: partRepoModel.Manufacturer.Website,
		},
		StockQuantity: partRepoModel.StockQuantity,
		Tags:          partRepoModel.Tags,
		CreatedAt:     partRepoModel.CreatedAt,
		UpdatedAt:     partRepoModel.UpdatedAt,
		Metadata:      &modelMetadata,
	}
}
