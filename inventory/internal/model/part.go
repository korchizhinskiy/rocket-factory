package model

import (
	"time"

	"github.com/google/uuid"
)

type Category int

const (
	CategoryUnknown Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimension struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}
type MetadataValue struct {
	StringValue *string
	IntValue    *int
	FloatValue  *float64
	BoolValue   *bool
}

type Part struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimension
	Manufacturer  Manufacturer
	Tags          *[]string
	Metadata      *map[string]MetadataValue
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
