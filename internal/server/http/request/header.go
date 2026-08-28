package request

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
)

// GetUUIDFromHeader extracts UUID from request header.
func GetUUIDFromHeader(r *http.Request, headerName string) (*uuid.UUID, error) {
	headerValue := r.Header.Get(headerName)
	if headerValue == "" {
		return &uuid.Nil, errs.NotFound("header not found")
	}

	userID, err := uuid.Parse(headerValue)
	if err != nil {
		return &uuid.Nil, errs.InvalidRequest(fmt.Sprintf("cannot parse uuid = '%s'", headerValue))
	}

	return &userID, nil
}
