package address

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
func (a Address) Value() (driver.Value, error) {
	if a.City == "" || a.Street == "" || a.House == "" {
		return nil, fmt.Errorf("invalid address: city, street and house are required")
	}

	b, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal address: %w", err)
	}

	return b, nil
}
