package request

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (string, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return "", errs.InvalidArgument(fmt.Sprintf("no key='%s' in path values", key))
	}

	if _, err := uuid.Parse(pathValue); err != nil {
		return "", errs.InvalidArgument(fmt.Sprintf("failed to parse path value='%s' by key='%s' to uuid: %s", pathValue, key, err))
	}

	return pathValue, nil
}
