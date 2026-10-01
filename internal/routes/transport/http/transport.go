package routes_transport_http

import (
	"context"
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
)

type RoutesHTTPHandler struct {
	routesService RoutesService
}

type RoutesService interface {
	GetBestRoutes(
		ctx context.Context,
		FromCityID uuid.UUID,
		ToCityID uuid.UUID,
		Sort *domain.SortType,
	) ([]domain.Route, error)
}

func NewRoutesHTTPHandler(
	routesService RoutesService,
) *RoutesHTTPHandler {
	return &RoutesHTTPHandler{
		routesService: routesService,
	}
}
