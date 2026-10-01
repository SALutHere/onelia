package core_http_request

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SALutHere/onelia/internal/core/domain"
	core_errors "github.com/SALutHere/onelia/internal/core/errors"
)

func GetRequiredStringQueryParam(r *http.Request, key string) (string, error) {
	param := strings.TrimSpace(r.URL.Query().Get(key))
	if param == "" {
		return "", fmt.Errorf(
			"param by key='%s' is required: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	return param, nil
}

func GetSortTypeQueryParam(r *http.Request, key string) *domain.SortType {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}

	val := domain.SortType(param)

	return &val
}
