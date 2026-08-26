package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
)

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %w: %w", err, errs.ErrInvalidRequest)
	}
	return nil
}
