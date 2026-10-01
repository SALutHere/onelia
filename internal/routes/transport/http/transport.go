package routes_transport_http

import (
	"context"

	"github.com/SALutHere/onelia/internal/core/domain"
)

type RoutesHTTPHandler struct {
	routesService RoutesService
}

type RoutesService interface {
	GetBestRoutes(
		ctx context.Context,
		fromCity, toCity string,
		order *domain.SortType,
	) ([]domain.Route, error)
}

func NewRoutesHTTPHandler(
	routesService RoutesService,
) *RoutesHTTPHandler {
	return &RoutesHTTPHandler{
		routesService: routesService,
	}
}
