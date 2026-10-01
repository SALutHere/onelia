package routes_postgres_repository

import (
	"context"
	"fmt"

	"github.com/SALutHere/onelia/internal/core/domain"
)

func (r *RoutesRepository) GetCities(ctx context.Context) ([]domain.City, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT
		id,
		name
	FROM onelia.cities
	ORDER BY id;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("load cities: %w", err)
	}
	defer rows.Close()

	var cityModels []CityModel
	for rows.Next() {
		var cityModel CityModel

		err := rows.Scan(
			&cityModel.ID,
			&cityModel.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("scan cities: %w", err)
		}

		cityModels = append(cityModels, cityModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	cityDomains := cityDomainsFromModels(cityModels)

	return cityDomains, nil
}
