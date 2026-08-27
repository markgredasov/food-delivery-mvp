package address

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
)

const maxLen = 100

// Address represents a delivery address for an order.
type Address struct {
	City      string `json:"city"`
	Street    string `json:"street"`
	House     string `json:"house"`
	Apartment string `json:"apartment"`
	Comment   string `json:"comment"`
}

// New validates and builds an Address.
func New(city, street, house, apartment string, comment *string) (Address, error) {
	cityTrimmed := strings.TrimSpace(city)
	if cityTrimmed == "" {
		return Address{}, errs.InvalidArgument("city must not be empty")
	}
	if utf8.RuneCountInString(cityTrimmed) > maxLen {
		return Address{}, errs.InvalidArgument("city is too long")
	}

	streetTrimmed := strings.TrimSpace(street)
	if streetTrimmed == "" {
		return Address{}, errs.InvalidArgument("street must not be empty")
	}
	if utf8.RuneCountInString(streetTrimmed) > maxLen {
		return Address{}, errs.InvalidArgument("street is too long")
	}

	houseTrimmed := strings.TrimSpace(house)
	if houseTrimmed == "" {
		return Address{}, errs.InvalidArgument("house must not be empty")
	}
	if utf8.RuneCountInString(houseTrimmed) > maxLen {
		return Address{}, errs.InvalidArgument("house is too long")
	}

	apartmentTrimmed := strings.TrimSpace(apartment)
	if apartmentTrimmed == "" {
		return Address{}, errs.InvalidArgument("apartment must not be empty")
	}
	if utf8.RuneCountInString(apartmentTrimmed) > maxLen {
		return Address{}, errs.InvalidArgument("apartment is too long")
	}

	address := Address{
		City:      cityTrimmed,
		Street:    streetTrimmed,
		House:     houseTrimmed,
		Apartment: apartmentTrimmed,
	}

	if comment != nil {
		commentTrimmed := strings.TrimSpace(*comment)
		if commentTrimmed == "" {
			return Address{}, errs.InvalidArgument("comment must not be empty")
		}
		if utf8.RuneCountInString(commentTrimmed) > maxLen {
			return Address{}, errs.InvalidArgument("comment is too long")
		}
		address.Comment = commentTrimmed
	}

	return address, nil
}

// Scan implements the [sql.Scanner] interface for JSONB data.
// It converts JSONB data from PostgreSQL into an Address struct.
func (a *Address) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("cannot scan nil into Address")
	}

	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into Address", src)
	}

	if err := json.Unmarshal(data, a); err != nil {
		return fmt.Errorf("failed to unmarshal address: %w", err)
	}

	return nil
}

// Value implements the [driver.Valuer] interface for JSONB storage.
// It serializes the Address struct to JSON for database storage.
func (a *Address) Value() (driver.Value, error) {
	if a.City == "" || a.Street == "" || a.House == "" {
		return nil, fmt.Errorf("invalid address: city, street and house are required")
	}

	b, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal address: %w", err)
	}

	return b, nil
}
