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
			if slices.Contains(part.Tags, filterTag) {
				filteredMap[partUUID] = part
			}
		}
	}
	return filteredMap
}
