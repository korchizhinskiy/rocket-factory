package filters

import (
	"slices"

	"github.com/google/uuid"

	"github.com/korchizhinskiy/rocket-factory/inventory/internal/model"
)

type PartUUIDFilter struct{}

func (f PartUUIDFilter) Filter(parts []model.Part, values []uuid.UUID) []model.Part {
	if len(values) == 0 {
		return parts
	}
	setUUID := make(map[uuid.UUID]struct{}, len(values))
	for _, v := range values {
		setUUID[v] = struct{}{}
	}

	filteredParts := make([]model.Part, 0)
	for _, p := range parts {
		if _, ok := setUUID[p.ID]; ok {
			filteredParts = append(filteredParts, p)
		}
	}
	return filteredParts
}

type PartTagFilter struct{}

func (f PartTagFilter) Filter(
	parts []model.Part,
	values []string,
) []model.Part {
	if len(values) == 0 {
		return parts
	}

	filteredParts := make([]model.Part, 0)

	for _, filterTag := range values {
		for _, part := range parts {
			if part.Tags != nil && slices.Contains(*part.Tags, filterTag) {
				filteredParts = append(filteredParts, part)
			}
		}
	}
	return filteredParts
}

type PartNameFilter struct{}

func (f PartNameFilter) Filter(
	parts []model.Part,
	values []string,
) []model.Part {
	if len(values) == 0 {
		return parts
	}

	filteredParts := make([]model.Part, 0)

	for _, filterName := range values {
		for _, part := range parts {
			if part.Name == filterName {
				filteredParts = append(filteredParts, part)
			}
		}
	}
	return filteredParts
}

type PartCategoryFilter struct{}

func (f PartCategoryFilter) Filter(
	parts []model.Part,
	values []string,
) []model.Part {
	if len(values) == 0 {
		return parts
	}

	filteredParts := make([]model.Part, 0)

	for _, filterCategory := range values {
		for _, part := range parts {
			if part.Category == filterCategory {
				filteredParts = append(filteredParts, part)
			}
		}
	}
	return filteredParts
}

type PartManufacturerCountriesFilter struct{}

func (f PartManufacturerCountriesFilter) Filter(
	parts []model.Part,
	values []string,
) []model.Part {
	if len(values) == 0 {
		return parts
	}

	filteredParts := make([]model.Part, 0)

	for _, filterManufacturerCountry := range values {
		for _, part := range parts {
			if part.Manufacturer.Country == filterManufacturerCountry {
				filteredParts = append(filteredParts, part)
			}
		}
	}
	return filteredParts
}
