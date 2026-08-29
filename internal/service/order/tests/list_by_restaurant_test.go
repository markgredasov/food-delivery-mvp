package orderservice_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/address"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
	service "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/service/order"
	mocks "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/mocks/order"
)

func TestListRestaurantOrders_Success(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)

	restaurantID := uuid.New()
	order1 := createTestOrder(uuid.New(), restaurantID, order.StatusPending)
	order2 := createTestOrder(uuid.New(), restaurantID, order.StatusAccepted)
	expectedOrders := []order.Order{order1, order2}

	rest := createTestRestaurant(restaurantID, restaurant.RestaurantStatusActive, "https://webhook.com")

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(rest, nil)

	orderMock.On("ListActiveByRestaurant", ctx, restaurantID).Return(expectedOrders, nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.ListRestaurantOrders(ctx, restaurantID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedOrders, result)

	restMock.AssertExpectations(t)
	orderMock.AssertExpectations(t)
}

func TestListRestaurantOrders_RestaurantNotFound(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)

	restaurantID := uuid.New()

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(restaurant.Restaurant{}, errs.NotFound("restaurant not found"))

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.ListRestaurantOrders(ctx, restaurantID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrNotFound)
	assert.Contains(t, err.Error(), "restaurant not found")
	assert.Nil(t, result)

	restMock.AssertExpectations(t)
	orderMock.AssertExpectations(t)
}

func TestListRestaurantOrders_RestaurantInactive(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)

	restaurantID := uuid.New()
	rest := createTestRestaurant(restaurantID, restaurant.RestaurantStatusInactive, "")

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(rest, nil)

	orderMock.On("ListActiveByRestaurant", ctx, restaurantID).Return([]order.Order{}, nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.ListRestaurantOrders(ctx, restaurantID)

	require.NoError(t, err)
	assert.Empty(t, result)

	restMock.AssertExpectations(t)
	orderMock.AssertExpectations(t)
}

func TestListRestaurantOrders_NoOrders(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)

	restaurantID := uuid.New()
	rest := createTestRestaurant(restaurantID, restaurant.RestaurantStatusActive, "https://webhook.com")

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(rest, nil)
	orderMock.On("ListActiveByRestaurant", ctx, restaurantID).Return([]order.Order{}, nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.ListRestaurantOrders(ctx, restaurantID)

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.Empty(t, result)

	restMock.AssertExpectations(t)
	orderMock.AssertExpectations(t)
}

func TestListRestaurantOrders_RepositoryError(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)

	restaurantID := uuid.New()
	rest := createTestRestaurant(restaurantID, restaurant.RestaurantStatusActive, "https://webhook.com")

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(rest, nil)

	orderMock.On("ListActiveByRestaurant", ctx, restaurantID).Return(nil, errs.Internal("database error", nil))

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.ListRestaurantOrders(ctx, restaurantID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInternal)
	assert.Contains(t, err.Error(), "database error")
	assert.Nil(t, result)

	restMock.AssertExpectations(t)
	orderMock.AssertExpectations(t)
}

func createTestRestaurant(id uuid.UUID, status restaurant.Status, serviceURL string) restaurant.Restaurant {
	addr, _ := address.New("Moscow", "Tverskaya", "1", "10", nil)
	rest, _ := restaurant.New(
		id,
		"Pizza House",
		nil,
		addr,
		30,
		status,
		serviceURL,
	)
	return rest
}
