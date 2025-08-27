package request

import (
	"github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
)

func ConvertPartModelToGetPartResponse(part model.Part) *inventoryv1.Part {
	responseMetadata := make(map[string]*inventoryv1.MetadataValue)

	if part.Metadata != nil {
		responseMetadata = make(map[string]*inventoryv1.MetadataValue, len(*part.Metadata))

		for k, v := range *part.Metadata {
			var mv *inventoryv1.MetadataValue

			switch {
			case v.StringValue != nil:
				mv = &inventoryv1.MetadataValue{
					Type: &inventoryv1.MetadataValue_StringValue{StringValue: *v.StringValue},
				}
			case v.IntValue != nil:
				mv = &inventoryv1.MetadataValue{
					Type: &inventoryv1.MetadataValue_Int64Value{Int64Value: int64(*v.IntValue)},
				}
			case v.FloatValue != nil:
				mv = &inventoryv1.MetadataValue{
					Type: &inventoryv1.MetadataValue_DoubleValue{DoubleValue: *v.FloatValue},
				}
			case v.BoolValue != nil:
				mv = &inventoryv1.MetadataValue{
					Type: &inventoryv1.MetadataValue_BoolValue{BoolValue: *v.BoolValue},
				}
			default:
				mv = &inventoryv1.MetadataValue{}
				continue
			}

			responseMetadata[k] = mv
		}
	}

	responseCategory := inventoryv1.PartCategory_PART_CATEGORY_UNSPECIFIED
	switch part.Category {
	case "PART_CATEGORY_ENGINE":
		responseCategory = inventoryv1.PartCategory_PART_CATEGORY_ENGINE
	case "PART_CATEGORY_FUEL":
		responseCategory = inventoryv1.PartCategory_PART_CATEGORY_FUEL
	case "PART_CATEGORY_PORTHOLE":
		responseCategory = inventoryv1.PartCategory_PART_CATEGORY_PORTHOLE
	case "PART_CATEGORY_WING":
		responseCategory = inventoryv1.PartCategory_PART_CATEGORY_WING
	}

	return &inventoryv1.Part{
		Uuid:        part.ID.String(),
		Name:        part.Name,
		Description: part.Description,
		Price:       part.Price,
		Dimensions: &inventoryv1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Category: responseCategory,
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:          *part.Tags,
		StockQuantity: part.StockQuantity,
		Metadata:      responseMetadata,
	}
}
