package catalogservice_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
	service "github.com/markgredasov/food-delivery-mvp/internal/service/catalog"
	mocks "github.com/markgredasov/food-delivery-mvp/mocks/catalog"
)

func TestGetMenu_Success(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	restaurantID := uuid.New()

	rest := createTestRestaurant(restaurantID, "https://webhook.com")

	menuItems := []menu.MenuItem{
		createTestMenuItem(uuid.New(), restaurantID, "Pizza", money.MustFromString("10.00"), true),
		createTestMenuItem(uuid.New(), restaurantID, "Pasta", money.MustFromString("12.50"), true),
		createTestMenuItem(uuid.New(), restaurantID, "Salad", money.MustFromString("8.00"), false),
	}

	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	categoryMock := mocks.NewMockCategoryRepository(t)
	txMock := mocks.NewMockTxManager(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(rest, nil)

	menuMock.On("ListByRestaurantID", ctx, restaurantID.String()).Return(menuItems, nil)

	svc := service.New(restMock, menuMock, categoryMock, txMock)

	result, err := svc.GetMenu(ctx, restaurantID)

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, menuItems, result)

	restMock.AssertExpectations(t)
	menuMock.AssertExpectations(t)
}

func TestGetMenu_RestaurantNotFound(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	restaurantID := uuid.New()

	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	categoryMock := mocks.NewMockCategoryRepository(t)
	txMock := mocks.NewMockTxManager(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(restaurant.Restaurant{}, errs.NotFound("restaurant not found"))

	svc := service.New(restMock, menuMock, categoryMock, txMock)

	result, err := svc.GetMenu(ctx, restaurantID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrNotFound)
	assert.Contains(t, err.Error(), "restaurant not found")
	assert.Nil(t, result)

	restMock.AssertExpectations(t)
	menuMock.AssertExpectations(t)
}

func TestGetMenu_EmptyMenu(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	restaurantID := uuid.New()

	rest := createTestRestaurant(restaurantID, "https://webhook.com")

	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	categoryMock := mocks.NewMockCategoryRepository(t)
	txMock := mocks.NewMockTxManager(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(rest, nil)
	menuMock.On("ListByRestaurantID", ctx, restaurantID.String()).Return([]menu.MenuItem{}, nil)

	svc := service.New(restMock, menuMock, categoryMock, txMock)

	result, err := svc.GetMenu(ctx, restaurantID)

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.Empty(t, result)

	restMock.AssertExpectations(t)
	menuMock.AssertExpectations(t)
}

func TestGetMenu_RepositoryError(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	restaurantID := uuid.New()

	rest := createTestRestaurant(restaurantID, "https://webhook.com")

	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	categoryMock := mocks.NewMockCategoryRepository(t)
	txMock := mocks.NewMockTxManager(t)

	restMock.On("GetByID", ctx, restaurantID.String()).Return(rest, nil)
	menuMock.On("ListByRestaurantID", ctx, restaurantID.String()).Return(nil, errs.Internal("database error", nil))

	svc := service.New(restMock, menuMock, categoryMock, txMock)

	result, err := svc.GetMenu(ctx, restaurantID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInternal)
	assert.Contains(t, err.Error(), "database error")
	assert.Nil(t, result)

	restMock.AssertExpectations(t)
	menuMock.AssertExpectations(t)
}
