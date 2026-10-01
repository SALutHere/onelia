package routes_postgres_repository

import (
	"context"
	"fmt"

	"github.com/SALutHere/onelia/internal/core/domain"
)

func (r *RoutesRepository) GetSegments(ctx context.Context) ([]domain.Segment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT
		id,
		from_city_id,
		to_city_id,
		transport_type,
		duration_minutes,
		(price * 100)::BIGINT
	FROM onelia.segments
	ORDER BY id;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("load segments: %w", err)
	}
	defer rows.Close()

	var segmentModels []SegmentModel
	for rows.Next() {
		var segmentModel SegmentModel

		err := rows.Scan(
			&segmentModel.ID,
			&segmentModel.FromCityID,
			&segmentModel.ToCityID,
			&segmentModel.TransportType,
			&segmentModel.DurationMinutes,
			&segmentModel.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("scan segments: %w", err)
		}

		segmentModels = append(segmentModels, segmentModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	segmentDomains := segmentDomainsFromModels(segmentModels)

	return segmentDomains, nil
}
