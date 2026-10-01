package routes_postgres_repository

import (
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
)

type CityModel struct {
	ID   uuid.UUID
	Name string
}

func cityDomainsFromModels(cities []CityModel) []domain.City {
	cityDomains := make([]domain.City, len(cities))

	for i, city := range cities {
		cityDomains[i] = domain.NewCity(
			city.ID,
			city.Name,
		)
	}

	return cityDomains
}

type SegmentModel struct {
	ID              uuid.UUID
	FromCityID      uuid.UUID
	ToCityID        uuid.UUID
	TransportType   string
	DurationMinutes int
	Price           int64
}

func segmentDomainsFromModels(segments []SegmentModel) []domain.Segment {
	segmentDomains := make([]domain.Segment, len(segments))

	for i, segment := range segments {
		segmentDomains[i] = domain.NewSegment(
			segment.ID,
			segment.FromCityID,
			segment.ToCityID,
			segment.TransportType,
			segment.DurationMinutes,
			segment.Price,
		)
	}

	return segmentDomains
}
