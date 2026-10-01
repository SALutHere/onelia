package routes_transport_http_dto

import "github.com/SALutHere/onelia/internal/core/domain"

type RouteDTOResponse struct {
	Parts                []RoutePartDTOResponse `json:"parts"`
	TotalDurationMinutes int                    `json:"total_duration_minutes" example:"560"`
	TotalPrice           int64                  `json:"total_price"            example:"550000"`
}

func RouteDTOFromDomain(route domain.Route) RouteDTOResponse {
	return RouteDTOResponse{
		Parts:                RoutePartsDTOFromDomains(route.Parts),
		TotalDurationMinutes: route.TotalDurationMinutes,
		TotalPrice:           route.TotalPrice,
	}
}

func RoutesDTOFromDomains(routes []domain.Route) []RouteDTOResponse {
	routesDTO := make([]RouteDTOResponse, len(routes))

	for i, route := range routes {
		routesDTO[i] = RouteDTOFromDomain(route)
	}

	return routesDTO
}
