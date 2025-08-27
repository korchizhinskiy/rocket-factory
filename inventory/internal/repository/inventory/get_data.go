package inventory

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"

	repoModel "github.com/korchizhinskiy/rocket-factory/inventory/internal/repository/model"
)

func ptr[T any](v T) *T { return &v }

func GenerateParts(n int) map[string]repoModel.Part {
	r := rand.New(rand.NewSource(42))
	categories := []string{"PART_CATEGORY_ENGINE", "PART_CATEGORY_FUEL", "PART_CATEGORY_PORTHOLE", "PART_CATEGORY_WING"}
	countries := []string{"RU", "CN", "TR", "DE"}

	m := make(map[string]repoModel.Part, n)
	now := time.Now().UTC()

	for i := 1; i <= n; i++ {
		id := uuid.NewString()
		name := fmt.Sprintf("Деталь %d", i)
		cat := categories[r.Intn(len(categories))]
		manName := fmt.Sprintf("Vendor%d", r.Intn(900)+100)
		country := countries[r.Intn(len(countries))]
		price := float64(r.Intn(2000)+50) / 100.0
		stock := int64(r.Intn(5000))

		dim := repoModel.Dimension{
			Length: 0.5 + r.Float64()*100,
			Width:  0.5 + r.Float64()*100,
			Height: 0.5 + r.Float64()*100,
			Weight: 0.05 + r.Float64()*10,
		}

		tags := []string{strings.ToLower(cat), "demo"}
		meta := map[string]repoModel.MetadataValue{
			"rating":   {FloatValue: ptr(r.Float64() * 5)},
			"featured": {BoolValue: ptr(r.Intn(2) == 0)},
		}

		m[id] = repoModel.Part{
			ID:            id,
			Name:          name,
			Description:   "Сгенерированная тестовая позиция",
			Price:         price,
			StockQuantity: stock,
			Category:      cat,
			Dimensions:    dim,
			Manufacturer: repoModel.Manufacturer{
				Name:    manName,
				Country: country,
				Website: fmt.Sprintf("https://%s.example", strings.ToLower(manName)),
			},
			Tags:      &tags,
			Metadata:  &meta,
			CreatedAt: now.Add(-time.Duration(r.Intn(180)) * 24 * time.Hour),
			UpdatedAt: now.Add(-time.Duration(r.Intn(30)) * 24 * time.Hour),
		}
	}
	return m
}
