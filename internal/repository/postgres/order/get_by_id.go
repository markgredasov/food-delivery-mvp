package orderrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	database_postgres "github.com/markgredasov/food-delivery-mvp/internal/database/postgres"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
)

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (order.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	q := r.pool.GetQuerier(ctx)

	sqlQuery := `
		SELECT id, restaurant_id, user_id, delivery_address, status, total_amount, comment, created_at, updated_at
		FROM kitchen.orders
		WHERE id = $1;
	`

	row := q.QueryRow(ctx, sqlQuery, id)
	o, err := scanOrder(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return order.Order{}, errs.NotFound(fmt.Sprintf("order with id = '%s'not found", id.String()))
		}
		return order.Order{}, errs.Internal("get order", err)
	}

	items, err := r.itemsFor(ctx, q, id)
	if err != nil {
		return order.Order{}, err
	}
	o.Items = items
	return o, nil
}

func (r *repository) itemsFor(ctx context.Context, q database_postgres.Querier, orderID uuid.UUID) ([]order.OrderItem, error) {
	sqlQuery := `
		SELECT oi.id, oi.order_id, oi.menu_item_id, mi.name, oi.quantity, oi.price
		FROM kitchen.order_items oi
		JOIN kitchen.menu_items mi ON mi.id = oi.menu_item_id
		WHERE oi.order_id = $1
		ORDER BY oi.id
	`
	rows, err := q.Query(ctx, sqlQuery, orderID)
	if err != nil {
		return nil, errs.Internal("list order items", err)
	}
	defer rows.Close()

	var out []order.OrderItem
	for rows.Next() {
		var it order.OrderItem
		var lineID uuid.UUID
		var price money.Money
		if err = rows.Scan(&lineID, &it.OrderID, &it.MenuItemID, &it.Name, &it.Quantity, &price); err != nil {
			return nil, errs.Internal("scan order item", err)
		}
		it.ID = lineID
		it.Price = price
		out = append(out, it)
	}
	if err = rows.Err(); err != nil {
		return nil, errs.Internal("iterate order items", err)
	}
	return out, nil
}
