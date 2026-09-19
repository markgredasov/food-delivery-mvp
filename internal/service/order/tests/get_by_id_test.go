package orderservice_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
	service "github.com/markgredasov/food-delivery-mvp/internal/service/order"
	mocks "github.com/markgredasov/food-delivery-mvp/mocks/order"
)

func TestGetOrder_Success(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restID := uuid.New()
	expectedOrder := createTestOrder(orderID, restID, order.StatusPending)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(expectedOrder, nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.GetOrder(ctx, orderID)

	require.NoError(t, err)
	assert.Equal(t, expectedOrder, result)
	assert.Equal(t, expectedOrder.ID, result.ID)
	assert.Equal(t, expectedOrder.RestaurantID, result.RestaurantID)
	assert.Equal(t, expectedOrder.Status, result.Status)

	orderMock.AssertExpectations(t)
}

func TestGetOrder_NotFound(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)

	orderID := uuid.New()

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(order.Order{}, errs.NotFound("order not found"))

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.GetOrder(ctx, orderID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrNotFound)
	assert.Contains(t, err.Error(), "order not found")
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}

func TestGetOrder_InternalError(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	dbError := errors.New("database connection failed")
	orderMock.On("GetByID", ctx, orderID).Return(order.Order{}, errs.Internal("get order", dbError))

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.GetOrder(ctx, orderID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInternal)
	assert.Contains(t, err.Error(), "get order")
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}

func TestGetOrder_EmptyID(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	emptyID := uuid.Nil

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, emptyID).Return(order.Order{}, errs.NotFound("order not found"))

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.GetOrder(ctx, emptyID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrNotFound)
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}
