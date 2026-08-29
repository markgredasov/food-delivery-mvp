package request

import (
	errs "avito-kitchen-restaurant-simulator/internal/errors"
	"encoding/json"
	"fmt"
	"net/http"
)

func Decode(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %w: %w", err, errs.ErrInvalidRequest)
	}
	return nil
}
