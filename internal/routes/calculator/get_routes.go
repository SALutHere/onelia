package routes_calculator

import (
	"context"
	"fmt"
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
	core_errors "github.com/SALutHere/onelia/internal/core/errors"
)

func (c *RoutesCalculator) GetRoutes(
	ctx context.Context,
	from, to uuid.UUID,
	maxSegments int,
) ([]domain.Route, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	for _, id := range []uuid.UUID{from, to} {
		if _, ok := c.cities[id]; !ok {
			return nil, fmt.Errorf("city %s: %w", id, core_errors.ErrNotFound)
		}
	}

	routes := make([]domain.Route, 0)
	visited := map[uuid.UUID]bool{from: true}
	path := make([]domain.RoutePart, 0)
	var walk func(uuid.UUID, int, int64) error
	walk = func(city uuid.UUID, duration int, price int64) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		if city == to && len(path) > 0 {
			routes = append(routes, domain.Route{
				Parts:                append([]domain.RoutePart(nil), path...),
				TotalDurationMinutes: duration,
				TotalPrice:           price,
			})
			return nil
		}
		if len(path) >= maxSegments {
			return nil
		}
		for _, segment := range c.outgoing[city] {
			if visited[segment.ToCityID] {
				continue
			}
			visited[segment.ToCityID] = true
			path = append(path, domain.RoutePart{
				Segment:  segment,
				FromCity: c.cities[city],
				ToCity:   c.cities[segment.ToCityID],
			})
			if err := walk(
				segment.ToCityID,
				duration+segment.DurationMinutes,
				price+segment.Price,
			); err != nil {
				return nil
			}
			path = path[:len(path)-1]
			delete(visited, segment.ToCityID)
		}
		return nil
	}

	if err := walk(from, 0, 0); err != nil {
		return nil, err
	}

	return routes, nil
}
