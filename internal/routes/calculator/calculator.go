package routes_calculator

import (
	"context"
	"fmt"
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
)

type RoutesCalculator struct {
	cities       map[uuid.UUID]domain.City
	citiesByName map[string]uuid.UUID
	outgoing     map[uuid.UUID][]domain.Segment
}

type RoutesRepository interface {
	GetCities(ctx context.Context) ([]domain.City, error)
	GetSegments(ctx context.Context) ([]domain.Segment, error)
}

func NewRoutesCalculator(
	ctx context.Context,
	routesRepository RoutesRepository,
) (*RoutesCalculator, error) {
	cities, err := routesRepository.GetCities(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cities from repository: %w", err)
	}

	segments, err := routesRepository.GetSegments(ctx)
	if err != nil {
		return nil, fmt.Errorf("get segments from repository: %w", err)
	}

	rc := &RoutesCalculator{
		cities:       make(map[uuid.UUID]domain.City, len(cities)),
		citiesByName: make(map[string]uuid.UUID, len(cities)),
		outgoing:     make(map[uuid.UUID][]domain.Segment),
	}
	for _, city := range cities {
		rc.cities[city.ID] = city
		rc.citiesByName[city.Name] = city.ID
	}
	for _, segment := range segments {
		if _, ok := rc.cities[segment.FromCityID]; !ok {
			return nil, fmt.Errorf("segment %s references unknown origin city", segment.ID)
		}
		if _, ok := rc.cities[segment.ToCityID]; !ok {
			return nil, fmt.Errorf("segment %s references unknown destination city", segment.ID)
		}
		rc.outgoing[segment.FromCityID] = append(rc.outgoing[segment.FromCityID], segment)
	}

	return rc, nil
}
