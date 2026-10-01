package routes_service

import (
	"context"
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
)

type RoutesService struct {
	routesCalculator RoutesCalculator
}

type RoutesCalculator interface {
	GetRoutes(
		ctx context.Context,
		from, to uuid.UUID,
		maxSegments int,
	) ([]domain.Route, error)
}

func NewRoutesService(
	routesCalculator RoutesCalculator,
) *RoutesService {
	return &RoutesService{
		routesCalculator: routesCalculator,
	}
}
