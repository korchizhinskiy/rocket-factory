package main

import (
	"math"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
)

func generateParts() map[string]*inventoryv1.Part {
	names := []string{
		"Main Engine",
		"Reserve Engine",
		"Thruster",
		"Fuel Tank",
		"Left Wing",
		"Right Wing",
		"Window A",
		"Window B",
		"Control Module",
		"Stabilizer",
	}
	descriptions := []string{
		"Primary propulsion unit",
		"Backup propulsion unit",
		"Thruster for fine adjustments",
		"Main fuel tank",
		"Left aerodynamic wing",
		"Right aerodynamic wing",
		"Front viewing window",
		"Side viewing window",
		"Flight control module",
		"Stabilization fin",
	}

	parts := make(map[string]*inventoryv1.Part, 0)
	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		genUUID := uuid.NewString()
		parts[genUUID] = &inventoryv1.Part{
			Uuid:        genUUID,
			Name:        names[idx],
			Description: descriptions[idx],
			Price: roundTo(
				gofakeit.Float64Range(100, 10_000),
			),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      inventoryv1.PartCategory(gofakeit.Number(1, 4)), // nolint: gosec
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     timestamppb.Now(),
		}
	}
	return parts
}

func generateDimensions() *inventoryv1.Dimensions {
	return &inventoryv1.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() *inventoryv1.Manufacturer {
	return &inventoryv1.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for range gofakeit.Number(1, 10) {
		tags = append(tags, gofakeit.EmojiTag())
	}
	return tags
}

func generateMetadata() map[string]*inventoryv1.MetadataValue {
	metadata := make(map[string]*inventoryv1.MetadataValue)
	for range gofakeit.Number(1, 10) {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}
	return metadata
}

func generateMetadataValue() *inventoryv1.MetadataValue {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &inventoryv1.MetadataValue{
			Type: &inventoryv1.MetadataValue_StringValue{
				StringValue: gofakeit.Word(),
			},
		}

	case 1:
		return &inventoryv1.MetadataValue{
			Type: &inventoryv1.MetadataValue_Int64Value{
				Int64Value: int64(gofakeit.Number(1, 100)),
			},
		}

	case 2:
		return &inventoryv1.MetadataValue{
			Type: &inventoryv1.MetadataValue_DoubleValue{
				DoubleValue: roundTo(
					gofakeit.Float64Range(1, 100),
				),
			},
		}

	case 3:
		return &inventoryv1.MetadataValue{
			Type: &inventoryv1.MetadataValue_BoolValue{
				BoolValue: gofakeit.Bool(),
			},
		}

	default:
		return nil
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}
