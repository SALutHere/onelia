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
	from, to string,
	maxSegments int,
) ([]domain.Route, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	fromID, ok := c.citiesByName[from]
	if !ok {
		return nil, fmt.Errorf("city %q: %w", from, core_errors.ErrNotFound)
	}
	toID, ok := c.citiesByName[to]
	if !ok {
		return nil, fmt.Errorf("city %q: %w", to, core_errors.ErrNotFound)
	}

	routes := make([]domain.Route, 0)
	visited := map[uuid.UUID]bool{fromID: true}
	path := make([]domain.RoutePart, 0)
	var walk func(uuid.UUID, int, int64) error
	walk = func(city uuid.UUID, duration int, price int64) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		if city == toID && len(path) > 0 {
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
				return err
			}
			path = path[:len(path)-1]
			delete(visited, segment.ToCityID)
		}
		return nil
	}

	if err := walk(fromID, 0, 0); err != nil {
		return nil, err
	}

	return routes, nil
}
