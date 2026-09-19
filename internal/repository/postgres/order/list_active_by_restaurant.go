package orderrepo

import (
	"context"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
)

var activeStatuses = []string{"pending", "sent_to_restaurant", "accepted", "preparing", "ready", "in_delivery"}

func (r *repository) ListActiveByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]order.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT 
			o.id, o.restaurant_id, o.user_id, o.delivery_address, o.status, o.total_amount, o.created_at, o.updated_at,
			oi.id, oi.menu_item_id, mi.name, oi.quantity, oi.price
		FROM kitchen.orders o
		LEFT JOIN kitchen.order_items oi ON oi.order_id = o.id
		LEFT JOIN kitchen.menu_items mi ON mi.id = oi.menu_item_id
		WHERE o.restaurant_id = $1 AND o.status = ANY($2)
		ORDER BY o.created_at DESC, oi.id;
	`

	rows, err := r.pool.GetQuerier(ctx).Query(ctx, sqlQuery, restaurantID, activeStatuses)
	if err != nil {
		return nil, errs.Internal("list active by restaurant", err)
	}
	defer rows.Close()

	ordersMap := make(map[uuid.UUID]*order.Order)
	var ordersList []order.Order

	for rows.Next() {
		var o order.Order
		var item order.OrderItem
		var itemID, menuItemID *uuid.UUID
		var itemName *string
		var quantity *int
		var price *money.Money

		err = rows.Scan(
			&o.ID, &o.RestaurantID, &o.UserID, &o.DeliveryAddress, &o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt,
			&itemID, &menuItemID, &itemName, &quantity, &price,
		)
		if err != nil {
			return nil, errs.Internal("scan order with items", err)
		}

		existingOrder, exists := ordersMap[o.ID]
		if !exists {
			o.Items = make([]order.OrderItem, 0)
			ordersMap[o.ID] = &o
			ordersList = append(ordersList, o)
			existingOrder = &ordersList[len(ordersList)-1]
		}

		if itemID != nil {
			item.ID = *itemID
			item.MenuItemID = *menuItemID
			item.Name = *itemName
			item.Quantity = *quantity
			item.Price = *price
			item.OrderID = o.ID
			existingOrder.Items = append(existingOrder.Items, item)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, errs.Internal("iterate orders with items", err)
	}

	return ordersList, nil
}
