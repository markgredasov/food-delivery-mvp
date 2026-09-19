package orderservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/model/address"
	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
)

func createTestOrder(id, restID uuid.UUID, status order.Status) order.Order {
	addr, _ := address.New("Moscow", "Tverskaya", "1", "10", nil)
	items := []order.OrderItem{
		{MenuItemID: uuid.New(), Quantity: 1, Price: money.MustFromString("10.00")},
	}
	ord, _ := order.New(id, restID, nil, addr, items, nil)
	ord.Status = status
	return ord
}

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
