package routes_transport_http_dto

import (
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
)

type RoutePartDTOResponse struct {
	ID              uuid.UUID       `json:"id"               example:"bd6bc4ad-cfcb-41b6-9e2c-adf3cbe21838"`
	FromCity        CityDTOResponse `json:"from_city"`
	ToCity          CityDTOResponse `json:"to_city"`
	TransportType   string          `json:"transport_type"   example:"train"`
	DurationMinutes int             `json:"duration_minutes" example:"360"`
	Price           int64           `json:"price"            example:"250000"`
}

func RoutePartDTOFromDomain(part domain.RoutePart) RoutePartDTOResponse {
	return RoutePartDTOResponse{
		ID:              part.ID,
		FromCity:        CityDTOFromDomain(part.FromCity),
		ToCity:          CityDTOFromDomain(part.ToCity),
		TransportType:   part.TransportType,
		DurationMinutes: part.DurationMinutes,
		Price:           part.Price,
	}
}

func RoutePartsDTOFromDomains(parts []domain.RoutePart) []RoutePartDTOResponse {
	routePartsDTO := make([]RoutePartDTOResponse, len(parts))

	for i, part := range parts {
		routePartsDTO[i] = RoutePartDTOFromDomain(part)
	}

	return routePartsDTO
}
