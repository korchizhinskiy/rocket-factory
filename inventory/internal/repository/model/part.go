package model

import "time"

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

type Category int

const (
	CategoryUnknown Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Part struct {
	ID            string
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
