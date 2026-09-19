package orderrepo

import (
	"context"

	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
)

func (r *repository) Create(ctx context.Context, o order.Order) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	q := r.pool.GetQuerier(ctx)

	sqlQueryInsertOrder := `
		INSERT INTO kitchen.orders (id, restaurant_id, user_id, delivery_address, status, total_amount, comment)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := q.Exec(ctx, sqlQueryInsertOrder, o.ID, o.RestaurantID, o.UserID,
		o.DeliveryAddress, o.Status.String(), o.TotalAmount.Decimal(), o.Comment)
	if err != nil {
		return errs.Internal("insert order", err)
	}

	sqlQueryInsertOrderItem := `
		INSERT INTO kitchen.order_items (id, order_id, menu_item_id, quantity, price)
		VALUES ($1, $2, $3, $4, $5);
	`

	for _, it := range o.Items {
		_, err = q.Exec(ctx, sqlQueryInsertOrderItem, it.ID, o.ID, it.MenuItemID, it.Quantity, it.Price)
		if err != nil {
			return errs.Internal("insert order item", err)
		}
	}
	return nil
}
