package orderrepo

import (
	"context"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
)

func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, from, to order.Status) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	q := r.pool.GetQuerier(ctx)

	sqlQuery := `
		UPDATE kitchen.orders
		SET status = $1
		WHERE id = $2 AND status = $3;
	`

	tag, err := q.Exec(ctx, sqlQuery, to.String(), id, from.String())
	if err != nil {
		return errs.Internal("update order status", err)
	}
	if tag.RowsAffected() == 0 {
		return errs.Conflict("order status changed concurrently, please retry")
	}
	return nil
}
