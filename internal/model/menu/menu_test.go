package menu_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMenuItem_RejectsEmptyName(t *testing.T) {
	_, err := menu.NewMenuItem(uuid.New(), uuid.New(), nil, "", nil, money.MustFromString("10.00"), true)
	if !errors.Is(err, errs.ErrInvalidArgument) {
		t.Fatalf("expected invalid error, got %v", err)
	}
}

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		inputName string
		wantName  string
		wantErr   error
	}{
		{
			name:      "valid category",
			id:        uuid.New(),
			inputName: "Супы",
			wantName:  "Супы",
			wantErr:   nil,
		},
		{
			name:      "trims whitespace",
			id:        uuid.New(),
			inputName: "  Супы  ",
			wantName:  "Супы",
			wantErr:   nil,
		},
		{
			name:      "empty name",
			id:        uuid.New(),
			inputName: "",
			wantErr:   errs.ErrInvalidArgument,
		},
		{
			name:      "whitespace only",
			id:        uuid.New(),
			inputName: "   ",
			wantErr:   errs.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat, err := menu.NewCategory(tt.id, tt.inputName)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cat.ID != tt.id {
				t.Errorf("expected ID %v, got %v", tt.id, cat.ID)
			}
			if cat.Name != tt.wantName {
				t.Errorf("expected name %q, got %q", tt.wantName, cat.Name)
			}
		})
	}
}

func TestNewMenuItem_Valid(t *testing.T) {
	id := uuid.New()
	restaurantID := uuid.New()
	category := &menu.Category{ID: uuid.New(), Name: "Soups"}
	price := money.MustFromString("350.00")

	tests := []struct {
		name        string
		description *string
		available   bool
	}{
		{
			name:        "without description",
			description: nil,
			available:   true,
		},
		{
			name:        "with description",
			description: new(string),
			available:   true,
		},
		{
			name:        "unavailable item",
			description: nil,
			available:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := menu.NewMenuItem(
				id,
				restaurantID,
				category,
				"Soup",
				tt.description,
				price,
				tt.available,
			)

			require.NoError(t, err)
			assert.Equal(t, id, item.ID)
			assert.Equal(t, restaurantID, item.RestaurantID)
			assert.Equal(t, category, item.Category)
			assert.Equal(t, "Soup", item.Name)
			assert.Equal(t, tt.description, item.Description)
			assert.Equal(t, price.String(), item.Price.String())
			assert.Equal(t, tt.available, item.Available)
		})
	}
}

func TestNewMenuItem_EmptyName(t *testing.T) {
	id := uuid.New()
	restaurantID := uuid.New()
	category := &menu.Category{ID: uuid.New(), Name: "Soups"}
	price := money.MustFromString("350.00")

	item, err := menu.NewMenuItem(
		id,
		restaurantID,
		category,
		"",
		nil,
		price,
		true,
	)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInvalidArgument)
	assert.Contains(t, err.Error(), "name must not be empty")
	assert.Equal(t, menu.MenuItem{}, item)
}

func TestNewMenuItem_ZeroPrice(t *testing.T) {
	id := uuid.New()
	restaurantID := uuid.New()
	category := &menu.Category{ID: uuid.New(), Name: "Soups"}
	zeroPrice := money.MustFromString("0.00")

	item, err := menu.NewMenuItem(
		id,
		restaurantID,
		category,
		"Soup",
		nil,
		zeroPrice,
		true,
	)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInvalidArgument)
	assert.Contains(t, err.Error(), "price must be greater than zero")
	assert.Equal(t, menu.MenuItem{}, item)
}

func TestNewMenuItem_NilCategory(t *testing.T) {
	id := uuid.New()
	restaurantID := uuid.New()
	price := money.MustFromString("350.00")

	item, err := menu.NewMenuItem(
		id,
		restaurantID,
		nil,
		"Soup",
		nil,
		price,
		true,
	)

	require.NoError(t, err)
	assert.Equal(t, id, item.ID)
	assert.Equal(t, restaurantID, item.RestaurantID)
	assert.Nil(t, item.Category)
	assert.Equal(t, "Soup", item.Name)
	assert.Equal(t, price.String(), item.Price.String())
	assert.True(t, item.Available)
}

func TestNewMenuItem_WhiteSpaceName(t *testing.T) {
	id := uuid.New()
	restaurantID := uuid.New()
	category := &menu.Category{ID: uuid.New(), Name: "Soups"}
	price := money.MustFromString("350.00")

	item, err := menu.NewMenuItem(
		id,
		restaurantID,
		category,
		"  soup  ",
		nil,
		price,
		true,
	)

	require.NoError(t, err)
	assert.Equal(t, "soup", item.Name)
}

func TestEnsureAvailable(t *testing.T) {
	tests := []struct {
		name      string
		available bool
		wantErr   error
	}{
		{
			name:      "available item",
			available: true,
			wantErr:   nil,
		},
		{
			name:      "unavailable item",
			available: false,
			wantErr:   errs.ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := menu.NewMenuItem(uuid.New(), uuid.New(), nil, "Борщ", nil, money.MustFromString("250.00"), tt.available)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			err = item.EnsureAvailable()
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestBelongsTo(t *testing.T) {
	restaurantID := uuid.New()
	otherRestaurantID := uuid.New()

	tests := []struct {
		name            string
		itemRestaurant  uuid.UUID
		checkRestaurant uuid.UUID
		wantResult      bool
	}{
		{
			name:            "belongs to restaurant",
			itemRestaurant:  restaurantID,
			checkRestaurant: restaurantID,
			wantResult:      true,
		},
		{
			name:            "does not belong to restaurant",
			itemRestaurant:  restaurantID,
			checkRestaurant: otherRestaurantID,
			wantResult:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := menu.NewMenuItem(uuid.New(), tt.itemRestaurant, nil, "Борщ", nil, money.MustFromString("250.00"), true)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			result := item.BelongsTo(tt.checkRestaurant)
			if result != tt.wantResult {
				t.Errorf("expected %v, got %v", tt.wantResult, result)
			}
		})
	}
}
