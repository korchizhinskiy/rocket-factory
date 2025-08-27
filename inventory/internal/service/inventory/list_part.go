package inventory

import (
	"context"

	"github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
	"github.com/korchizhinskiy/rocket-factory/inventory/internal/service/contract"
	partFilters "github.com/korchizhinskiy/rocket-factory/inventory/internal/service/filters"
)

func (s *service) GetListPart(ctx context.Context, filters contract.ListPartFilter) ([]model.Part, error) {
	uuidFilter := partFilters.PartUUIDFilter{}
	tagFilter := partFilters.PartTagFilter{}
	nameFilter := partFilters.PartNameFilter{}
	// categoryFilter := partFilters.PartCategoryFilter{}
	// manufacturerCountriesFilter := partFilters.PartManufacturerCountriesFilter{}
	parts, _ := s.inventoryRepository.GetListPart(ctx, filters)

	if filters.UUIDs != nil {
		parts = uuidFilter.Filter(
			parts,
			*filters.UUIDs,
		)
	}
	if filters.Tags != nil {
		parts = tagFilter.Filter(
			parts,
			*filters.Tags,
		)
	}
	if filters.Names != nil {
		parts = nameFilter.Filter(
			parts,
			*filters.Names,
		)
	}

	if filters.Categories != nil {
		categoryFilter := partFilters.PartCategoryFilter{}
		parts = categoryFilter.Filter(
			parts,
			*filters.Categories,
		)
	}

	if filters.ManufactorerCountries != nil {
		manufacturerCountriesFilter := partFilters.PartManufacturerCountriesFilter{}
		parts = manufacturerCountriesFilter.Filter(
			parts,
			*filters.ManufactorerCountries,
		)
	}

	// parts = tagFilter.filter(
	// 	parts,
	// 	request.GetFilter().Tags,
	// )
	// parts = nameFilter.filter(
	// 	parts,
	// 	request.GetFilter().Names,
	// )
	// parts = categoryFilter.filter(
	// 	parts,
	// 	request.GetFilter().Categories,
	// )
	// parts = manufacturerCountriesFilter.filter(
	// 	parts,
	// 	request.GetFilter().ManufactorerCountries,
	// )

	return parts, nil
}
