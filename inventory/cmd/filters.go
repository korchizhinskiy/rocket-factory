package main

import (
	"slices"

	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
)

type partUUIDFilter struct{}

func (f partUUIDFilter) filter(
	parts map[string]*inventoryv1.Part,
	values []string,
) map[string]*inventoryv1.Part {
	if len(values) == 0 {
		return parts
	}

	filteredMap := make(map[string]*inventoryv1.Part, 0)
	for _, filterUUID := range values {
		part, ok := parts[filterUUID]
		if ok {
			filteredMap[filterUUID] = part
		}
	}
	return filteredMap
}

type partTagFilter struct{}

func (f partTagFilter) filter(
	parts map[string]*inventoryv1.Part,
	values []string,
) map[string]*inventoryv1.Part {
	if len(values) == 0 {
		return parts
	}

	filteredMap := make(map[string]*inventoryv1.Part, 0)

	for _, filterTag := range values {
		for partUUID, part := range parts {
			if _, ok := filteredMap[partUUID]; !ok && slices.Contains(part.Tags, filterTag) {
				filteredMap[partUUID] = part
			}
		}
	}
	return filteredMap
}

type partNameFilter struct{}

func (f partNameFilter) filter(
	parts map[string]*inventoryv1.Part,
	values []string,
) map[string]*inventoryv1.Part {
	if len(values) == 0 {
		return parts
	}

	filteredMap := make(map[string]*inventoryv1.Part, 0)

	for _, filterName := range values {
		for partUUID, part := range parts {
			if _, ok := filteredMap[partUUID]; !ok && part.Name == filterName {
				filteredMap[partUUID] = part
			}
		}
	}
	return filteredMap
}

type partCategoryFilter struct{}

func (f partCategoryFilter) filter(
	parts map[string]*inventoryv1.Part,
	values []inventoryv1.PartCategory,
) map[string]*inventoryv1.Part {
	if len(values) == 0 {
		return parts
	}

	filteredMap := make(map[string]*inventoryv1.Part, 0)

	for _, filterCategory := range values {
		for partUUID, part := range parts {
			if _, ok := filteredMap[partUUID]; !ok && part.Category == filterCategory {
				filteredMap[partUUID] = part
			}
		}
	}
	return filteredMap
}

type partManufacturerCountriesFilter struct{}

func (f partManufacturerCountriesFilter) filter(
	parts map[string]*inventoryv1.Part,
	values []string,
) map[string]*inventoryv1.Part {
	if len(values) == 0 {
		return parts
	}

	filteredMap := make(map[string]*inventoryv1.Part, 0)

	for _, filterManufacturerCountry := range values {
		for partUUID, part := range parts {
			if _, ok := filteredMap[partUUID]; !ok &&
				part.Manufacturer.Country == filterManufacturerCountry {
				filteredMap[partUUID] = part
			}
		}
	}
	return filteredMap
}
