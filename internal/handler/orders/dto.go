package ordershandler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/address"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/money"
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

type OrderItemDTO struct {
	ID         string      `json:"id"`
	OrderID    string      `json:"order_id"`
	MenuItemID string      `json:"menu_item_id"`
	Name       string      `json:"name"`
	Price      money.Money `json:"price"`
	Quantity   int         `json:"quantity"`
}

func toOrderItem(m order.OrderItem) OrderItemDTO {
	return OrderItemDTO{
		ID:         m.ID.String(),
		OrderID:    m.OrderID.String(),
		MenuItemID: m.MenuItemID.String(),
		Name:       m.Name,
		Price:      m.Price,
		Quantity:   m.Quantity,
	}
}

type OrderDTO struct {
	ID              string          `json:"id"`
	RestaurantID    string          `json:"restaurant_id"`
	UserID          *string         `json:"user_id"`
	Status          string          `json:"status"`
	DeliveryAddress address.Address `json:"delivery_address"`
	TotalAmount     money.Money     `json:"total_amount"`
	Items           []OrderItemDTO  `json:"order_items"`
	Comment         *string         `json:"comment"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func toOrder(m order.Order) OrderDTO {
	items := make([]OrderItemDTO, len(m.Items))
	for i := range m.Items {
		items[i] = toOrderItem(m.Items[i])
	}

	var uid *string
	if m.UserID != nil {
		userIDString := m.UserID.String()
		uid = &userIDString
	}

	return OrderDTO{
		ID:              m.ID.String(),
		RestaurantID:    m.RestaurantID.String(),
		UserID:          uid,
		Status:          m.Status.String(),
		DeliveryAddress: m.DeliveryAddress,
		TotalAmount:     m.TotalAmount,
		Items:           items,
		Comment:         m.Comment,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

type OrdersDTO struct {
	Orders []OrderDTO `json:"orders"`
}

func toOrders(m []order.Order) OrdersDTO {
	out := make([]OrderDTO, len(m))
	for i := range m {
		out[i] = toOrder(m[i])
	}

	return OrdersDTO{
		Orders: out,
	}
}

type RejectOrderDTO struct {
	Reason string `json:"reason"`
}
