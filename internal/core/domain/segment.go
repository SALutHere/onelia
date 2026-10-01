package domain

import "uuid"

type Segment struct {
	ID              uuid.UUID
	FromCityID      uuid.UUID
	ToCityID        uuid.UUID
	TransportType   string
	DurationMinutes int
	Price           int64 // в копейках
}

func NewSegment(
	id uuid.UUID,
	fromCityID uuid.UUID,
	toCityID uuid.UUID,
	transportType string,
	durationMinutes int,
	price int64,
) Segment {
	return Segment{
		ID:              id,
		FromCityID:      fromCityID,
		ToCityID:        toCityID,
		TransportType:   transportType,
		DurationMinutes: durationMinutes,
		Price:           price,
	}
}
