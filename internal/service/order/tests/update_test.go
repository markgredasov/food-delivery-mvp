package orderservice_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
	service "github.com/markgredasov/food-delivery-mvp/internal/service/order"
	mocks "github.com/markgredasov/food-delivery-mvp/mocks/order"
)

func TestAcceptOrder_Success(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()

	ord := createTestOrder(orderID, restaurantID, order.StatusPending)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(ord, nil)

	orderMock.On("UpdateStatus", ctx, orderID, order.StatusPending, order.StatusAccepted).Return(nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.AcceptOrder(ctx, restaurantID, orderID)

	require.NoError(t, err)
	assert.Equal(t, order.StatusAccepted, result.Status)
	assert.Equal(t, orderID, result.ID)

	orderMock.AssertExpectations(t)
}

func TestAcceptOrder_OrderNotFound(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(order.Order{}, errs.NotFound("order not found"))

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.AcceptOrder(ctx, restaurantID, orderID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrNotFound)
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}

func TestAcceptOrder_NotOwnedByRestaurant(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()
	otherRestaurantID := uuid.New()

	ord := createTestOrder(orderID, otherRestaurantID, order.StatusPending)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(ord, nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.AcceptOrder(ctx, restaurantID, orderID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrForbidden)
	assert.Contains(t, err.Error(), "order does not belong to this restaurant")
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}

func TestAcceptOrder_WrongStatus(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()

	ord := createTestOrder(orderID, restaurantID, order.StatusAccepted)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(ord, nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.AcceptOrder(ctx, restaurantID, orderID)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrConflict)
	assert.Contains(t, err.Error(), "cannot be accepted")
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}

func TestRejectOrder_Success(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()
	reason := "Out of stock"

	ord := createTestOrder(orderID, restaurantID, order.StatusPending)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(ord, nil)
	orderMock.On("UpdateStatus", ctx, orderID, order.StatusPending, order.StatusRejected).Return(nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.RejectOrder(ctx, restaurantID, orderID, reason)

	require.NoError(t, err)
	assert.Equal(t, order.StatusRejected, result.Status)
	assert.Equal(t, orderID, result.ID)

	orderMock.AssertExpectations(t)
}

func TestRejectOrder_NotOwnedByRestaurant(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()
	otherRestaurantID := uuid.New()
	reason := "Out of stock"

	ord := createTestOrder(orderID, otherRestaurantID, order.StatusPending)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(ord, nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.RejectOrder(ctx, restaurantID, orderID, reason)

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrForbidden)
	assert.Contains(t, err.Error(), "order does not belong to this restaurant")
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}

func TestUpdateStatus_Success(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()

	ord := createTestOrder(orderID, restaurantID, order.StatusAccepted)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(ord, nil)
	orderMock.On("UpdateStatus", ctx, orderID, order.StatusAccepted, order.StatusPreparing).Return(nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.UpdateStatus(ctx, restaurantID, orderID, "preparing")

	require.NoError(t, err)
	assert.Equal(t, order.StatusPreparing, result.Status)
	assert.Equal(t, orderID, result.ID)

	orderMock.AssertExpectations(t)
}

func TestUpdateStatus_InvalidStatus(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.UpdateStatus(ctx, restaurantID, orderID, "invalid_status")

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInvalidArgument)
	assert.Contains(t, err.Error(), "invalid status")
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}

func TestUpdateStatus_NotOwnedByRestaurant(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()
	otherRestaurantID := uuid.New()

	ord := createTestOrder(orderID, otherRestaurantID, order.StatusAccepted)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(ord, nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.UpdateStatus(ctx, restaurantID, orderID, "preparing")

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrForbidden)
	assert.Contains(t, err.Error(), "order does not belong to this restaurant")
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}

func TestUpdateStatus_WrongTransition(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()

	ord := createTestOrder(orderID, restaurantID, order.StatusPending)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(ord, nil)

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.UpdateStatus(ctx, restaurantID, orderID, "preparing")

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrConflict)
	assert.Contains(t, err.Error(), "cannot transition")
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}

func TestUpdateStatus_ConcurrentUpdate(t *testing.T) {
	t.Parallel()

	ctx := ctxWithTestLogger(t)
	orderID := uuid.New()
	restaurantID := uuid.New()

	ord := createTestOrder(orderID, restaurantID, order.StatusAccepted)

	txMock := mocks.NewMockTxManager(t)
	restMock := mocks.NewMockRestaurantRepository(t)
	menuMock := mocks.NewMockMenuRepository(t)
	orderMock := mocks.NewMockOrderRepository(t)
	webhookMock := mocks.NewMockWebhookSender(t)

	orderMock.On("GetByID", ctx, orderID).Return(ord, nil)
	orderMock.On("UpdateStatus", ctx, orderID, order.StatusAccepted, order.StatusPreparing).
		Return(errs.Conflict("order status changed concurrently, please retry"))

	svc := service.New(restMock, menuMock, orderMock, txMock, webhookMock)

	result, err := svc.UpdateStatus(ctx, restaurantID, orderID, "preparing")

	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrConflict)
	assert.Contains(t, err.Error(), "order status changed concurrently")
	assert.Equal(t, order.Order{}, result)

	orderMock.AssertExpectations(t)
}
