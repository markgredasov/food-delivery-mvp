package address

import (
	"strings"
	"unicode/utf8"

	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
)

const maxLen = 100

// Address represents a delivery address for an order.
type Address struct {
	City      string
	Street    string
	House     string
	Apartment string
	Comment   string
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
