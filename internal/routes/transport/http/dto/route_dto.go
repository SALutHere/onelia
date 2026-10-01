package routes_transport_http_dto

import "github.com/SALutHere/onelia/internal/core/domain"

type RouteDTOResponse struct {
	Parts                []RoutePartDTOResponse `json:"parts"`
	TotalDurationMinutes int                    `json:"total_duration_minutes"`
	TotalPrice           int64                  `json:"total_price"`
}

func RouteDTOFromDomain(route domain.Route) RouteDTOResponse {
	var totalDuration int
	var totalPrice int64

	for _, part := range route.Parts {
		totalDuration += part.DurationMinutes
		totalPrice += part.Price
	}

	return RouteDTOResponse{
		Parts:                RoutePartsDTOFromDomains(route.Parts),
		TotalDurationMinutes: totalDuration,
		TotalPrice:           totalPrice,
	}
}

func RoutesDTOFromDomains(routes []domain.Route) []RouteDTOResponse {
	routesDTO := make([]RouteDTOResponse, len(routes))

	for i, route := range routes {
		routesDTO[i] = RouteDTOFromDomain(route)
	}

	return routesDTO
}
