package contract

import "github.com/google/uuid"

type ListPartFilter struct {
	UUIDs                 *[]uuid.UUID
	Tags                  *[]string
	Names                 *[]string
	Categories            *[]string
	ManufactorerCountries *[]string
}
