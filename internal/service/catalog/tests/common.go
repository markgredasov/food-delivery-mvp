package catalogservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/address"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/money"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
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
