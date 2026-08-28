package ordershandler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/address"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
)

type AddressDTO struct {
	City      string  `json:"city"`
	Street    string  `json:"street"`
	House     string  `json:"house"`
	Apartment string  `json:"apartment"`
	Comment   *string `json:"comment"`
}

type CreateOrderItemDTORequest struct {
	MenuItemID string `json:"menu_item_id"`
	Quantity   int    `json:"quantity"`
}

type CreateOrderDTORequest struct {
	RestaurantID    string                      `json:"restaurant_id"`
	DeliveryAddress AddressDTO                  `json:"delivery_address"`
	Items           []CreateOrderItemDTORequest `json:"order_items"`
	Comment         *string                     `json:"comment"`
}

func createOrderToModel(req CreateOrderDTORequest, userID *uuid.UUID) (order.Order, error) {
	items := make([]order.OrderItem, len(req.Items))
	for i, reqItem := range req.Items {
		menuItemID, err := uuid.Parse(reqItem.MenuItemID)
		if err != nil {
			return order.Order{}, errs.InvalidRequest(fmt.Sprintf("cannot parse uuid = '%s'", reqItem.MenuItemID))
		}

		items[i], err = order.NewOrderItem(menuItemID, reqItem.Quantity)
		if err != nil {
			return order.Order{}, err
		}
	}

	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		return order.Order{}, errs.InvalidRequest(fmt.Sprintf("cannot parse uuid = '%s'", req.RestaurantID))
	}

	address, err := address.New(req.DeliveryAddress.City, req.DeliveryAddress.Street, req.DeliveryAddress.House,
		req.DeliveryAddress.Apartment, req.DeliveryAddress.Comment)
	if err != nil {
		return order.Order{}, err
	}

	return order.New(uuid.New(), restaurantID, userID, address, items, req.Comment)
}

type CreateOrderDTOResponse struct {
	ID                   string    `json:"id"`
	Status               string    `json:"status"`
	TotalAmount          string    `json:"total_amount"`
	EstimatedDevieryTime time.Time `json:"estimated_delivery_time"`
}

func createOrderToDTO(m order.Order) CreateOrderDTOResponse {
	return CreateOrderDTOResponse{
		ID:                   m.ID.String(),
		Status:               m.Status.String(),
		TotalAmount:          m.TotalAmount.String(),
		EstimatedDevieryTime: m.CreatedAt,
	}
}
