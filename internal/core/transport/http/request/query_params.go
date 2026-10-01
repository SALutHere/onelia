package core_http_request

import (
	"fmt"
	"net/http"
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
	core_errors "github.com/SALutHere/onelia/internal/core/errors"
)

func GetRequiredUUIDQueryParam(r *http.Request, key string) (uuid.UUID, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return uuid.Nil(), fmt.Errorf(
			"param by key='%s' is required: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	val, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil(), fmt.Errorf(
			"param='%s' by key='%s' not a valid UUID: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return val, nil
}

func GetSortTypeQueryParam(r *http.Request, key string) *domain.SortType {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}

	val := domain.SortType(param)

	return &val
}
