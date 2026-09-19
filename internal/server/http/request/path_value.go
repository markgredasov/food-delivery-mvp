package request

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil, errs.InvalidRequest(fmt.Sprintf("no key='%s' in path values", key))
	}

	id, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, errs.InvalidRequest(fmt.Sprintf("failed to parse path value='%s' by key='%s' to uuid: %s", pathValue, key, err))
	}

	return id, nil
}
