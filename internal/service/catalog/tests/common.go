package catalogservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/model/address"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
)

func ctxWithTestLogger(t *testing.T) context.Context {
	t.Helper()

	ctx := context.Background()
	tempDir := t.TempDir()

	cfg := logger.Config{
		Level:  "DEBUG",
		Folder: tempDir,
	}

	testLogger, err := logger.NewLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}

	t.Cleanup(func() {
		if err = testLogger.Close(); err != nil {
			t.Logf("failed to close logger: %v", err)
		}
	})

	return context.WithValue(ctx, "logger", testLogger) //nolint:staticcheck // not needed
}

func createTestRestaurant(id uuid.UUID, serviceURL string) restaurant.Restaurant {
	addr, _ := address.New("Moscow", "Tverskaya", "1", "10", nil)
	rest, _ := restaurant.New(
		id,
		"Pizza House",
		nil,
		addr,
		30,
		restaurant.RestaurantStatusActive,
		serviceURL,
	)
	return rest
}

func createTestMenuItem(id, restID uuid.UUID, name string, price money.Money, available bool) menu.MenuItem {
	item, _ := menu.NewMenuItem(
		id,
		restID,
		nil,
		name,
		nil,
		price,
		available,
	)
	return item
}
