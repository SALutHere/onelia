package routes_transport_http_dto

import (
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
)

type RoutePartDTOResponse struct {
	ID              uuid.UUID       `json:"id"`
	FromCity        CityDTOResponse `json:"from_city"`
	ToCity          CityDTOResponse `json:"to_city"`
	TransportType   string          `json:"transport_type"`
	DurationMinutes int             `json:"duration_minutes"`
	Price           int64           `json:"price"`
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
