package restaurant

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/address"
)

const maxNameLen = 100
const maxDescriptionLen = 500

// Status is the operating status of a restaurant.
type Status string

// Restaurant operating statuses; only StatusActive restaurants accept orders.
const (
	RestaurantStatusActive   Status = "active"
	RestaurantStatusInactive Status = "inactive"
	RestaurantStatusBlocked  Status = "blocked"
)

// NewStatus builds and validates RestaurantStatus
func NewStatus(s string) (Status, error) {
	status := Status(s)
	switch status {
	case RestaurantStatusActive, RestaurantStatusInactive, RestaurantStatusBlocked:
		return status, nil
	}

	return "", errs.InvalidArgument(fmt.Sprintf("status '%s'", s))
}

// String converts status to string.
func (s Status) String() string {
	return string(s)
}

// Restaurant is the aggregate root for a food establishment onboarded onto
// the platform.
type Restaurant struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Address     address.Address
	Status      Status
	ServiceURL  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// New validates and builds Restaurant.
func New(id uuid.UUID, name string, description *string, address address.Address,
	deliveryTime int, status Status, serviceURL string) (Restaurant, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return Restaurant{}, errs.InvalidArgument("name must not be empty")
	}
	if utf8.RuneCountInString(trimmedName) > maxNameLen {
		return Restaurant{}, errs.InvalidArgument("name is too long")
	}

	if deliveryTime <= 0 {
		return Restaurant{}, errs.InvalidArgument("delivery time must not be zero or less")
	}

	trimmedServiceURL := strings.TrimSpace(serviceURL)
	if trimmedServiceURL == "" {
		return Restaurant{}, errs.InvalidArgument("service url must not be empty")
	}

	restaurant := Restaurant{
		ID:         id,
		Name:       name,
		Address:    address,
		Status:     status,
		ServiceURL: serviceURL,
	}

	if description != nil {
		trimmedDescription := strings.TrimSpace(*description)
		if utf8.RuneCountInString(trimmedDescription) > maxDescriptionLen {
			return Restaurant{}, errs.InvalidArgument("description is too long")
		}
		restaurant.Description = &trimmedDescription
	}

	return restaurant, nil
}

// IsActive reports if restaurant is not blocked and not inactive.
func (r *Restaurant) IsActive() bool {
	return r.Status == RestaurantStatusActive
}

// EnsureAcceptsOrders returns a domain error if the restaurant cannot accept
// new orders right now.
func (r *Restaurant) EnsureAcceptsOrders() error {
	if !r.IsActive() {
		return errs.Conflict("restaurant is not accepting orders")
	}
	return nil
}

// HasServiceURL reports whether the restaurant configured a webhook URL for
// push order delivery.
func (r *Restaurant) HasServiceURL() bool { return r.ServiceURL != "" }
