package routes_service

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
	core_errors "github.com/SALutHere/onelia/internal/core/errors"
)

const (
	maxSegments = 3
	maxResults  = 3
)

func (s *RoutesService) GetBestRoutes(
	ctx context.Context,
	FromCityID uuid.UUID,
	ToCityID uuid.UUID,
	Sort *domain.SortType,
) ([]domain.Route, error) {
	if FromCityID == uuid.Nil() || ToCityID == uuid.Nil() {
		return nil, fmt.Errorf("both city IDs are required: %w", core_errors.ErrInvalidArgument)
	}
	if FromCityID == ToCityID {
		return nil, fmt.Errorf("origin and destination must be differ: %w", core_errors.ErrInvalidArgument)
	}
	if Sort == nil {
		Sort = new(domain.SortTime)
	}
	if *Sort != domain.SortTime && *Sort != domain.SortPrice {
		return nil, fmt.Errorf("sort must be 'time' or 'price': %w", core_errors.ErrInvalidArgument)
	}
	routes, err := s.routesCalculator.GetRoutes(ctx, FromCityID, ToCityID, maxSegments)
	if err != nil {
		return nil, fmt.Errorf("calculate routes: %w", err)
	}

	result := append(make([]domain.Route, 0, len(routes)), routes...)
	sort.Slice(result, func(i, j int) bool { return less(result[i], result[j], *Sort) })
	if len(result) > maxResults {
		result = result[:maxResults]
	}
	return result, nil
}

func less(a, b domain.Route, order domain.SortType) bool {
	if order == domain.SortPrice && a.TotalPrice != b.TotalPrice {
		return a.TotalPrice < b.TotalPrice
	}
	if a.TotalDurationMinutes != b.TotalDurationMinutes {
		return a.TotalDurationMinutes < b.TotalDurationMinutes
	}
	if a.TotalPrice != b.TotalPrice {
		return a.TotalPrice < b.TotalPrice
	}
	if len(a.Parts) != len(b.Parts) {
		return len(a.Parts) < len(b.Parts)
	}
	for i := range a.Parts {
		if cmp := bytes.Compare(a.Parts[i].ID[:], b.Parts[i].ID[:]); cmp != 0 {
			return cmp < 0
		}
	}
	return false
}
