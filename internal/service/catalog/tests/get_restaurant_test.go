package catalogservice_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
	service "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/service/catalog"
	mocks "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/mocks/catalog"
)

func TestGetRestaurant_Success(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	restaurantID := uuid.New()

	expectedRest := createTestRestaurant(restaurantID, "https://webhook.com")

	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	categoryMock := mocks.NewMockCategoryRepository(t)
	txMock := mocks.NewMockTxManager(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(expectedRest, nil)

	svc := service.New(restMock, menuMock, categoryMock, txMock)

	result, err := svc.GetRestaurant(ctx, restaurantID)

	require.NoError(t, err)
	assert.Equal(t, expectedRest, result)
	assert.Equal(t, expectedRest.ID, result.ID)
	assert.Equal(t, expectedRest.Name, result.Name)
	assert.Equal(t, expectedRest.Status, result.Status)

	restMock.AssertExpectations(t)
}

func TestGetRestaurant_NotFound(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	restaurantID := uuid.New()

	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	categoryMock := mocks.NewMockCategoryRepository(t)
	txMock := mocks.NewMockTxManager(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(restaurant.Restaurant{}, errs.NotFound("restaurant not found"))

	svc := service.New(restMock, menuMock, categoryMock, txMock)

	result, err := svc.GetRestaurant(ctx, restaurantID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrNotFound)
	assert.Contains(t, err.Error(), "restaurant not found")
	assert.Equal(t, restaurant.Restaurant{}, result)

	restMock.AssertExpectations(t)
}

func TestGetRestaurant_InternalError(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	restaurantID := uuid.New()

	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	categoryMock := mocks.NewMockCategoryRepository(t)
	txMock := mocks.NewMockTxManager(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(restaurant.Restaurant{}, errs.Internal("database connection failed", nil))

	svc := service.New(restMock, menuMock, categoryMock, txMock)

	result, err := svc.GetRestaurant(ctx, restaurantID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInternal)
	assert.Contains(t, err.Error(), "database connection failed")
	assert.Equal(t, restaurant.Restaurant{}, result)

	restMock.AssertExpectations(t)
}

func TestGetRestaurant_InvalidUUID(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	invalidID := uuid.Nil

	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	categoryMock := mocks.NewMockCategoryRepository(t)
	txMock := mocks.NewMockTxManager(t)

	restMock.On("GetByID", ctx, invalidID.String()).Return(restaurant.Restaurant{}, errs.InvalidArgument("invalid restaurant id"))

	svc := service.New(restMock, menuMock, categoryMock, txMock)

	result, err := svc.GetRestaurant(ctx, invalidID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInvalidArgument)
	assert.Contains(t, err.Error(), "invalid restaurant id")
	assert.Equal(t, restaurant.Restaurant{}, result)

	restMock.AssertExpectations(t)
}

func TestGetRestaurant_WithServiceURL(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	restaurantID := uuid.New()

	expectedRest := createTestRestaurant(restaurantID, "https://my-pizza.com/webhook")

	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	categoryMock := mocks.NewMockCategoryRepository(t)
	txMock := mocks.NewMockTxManager(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(expectedRest, nil)

	svc := service.New(restMock, menuMock, categoryMock, txMock)

	result, err := svc.GetRestaurant(ctx, restaurantID)

	require.NoError(t, err)
	assert.Equal(t, expectedRest, result)
	assert.Equal(t, "https://my-pizza.com/webhook", result.ServiceURL)
	assert.True(t, result.HasServiceURL())

	restMock.AssertExpectations(t)
}
