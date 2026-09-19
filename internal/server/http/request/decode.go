package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
)

func Decode(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %w: %w", err, errs.ErrInvalidRequest)
	}
	return nil
}
