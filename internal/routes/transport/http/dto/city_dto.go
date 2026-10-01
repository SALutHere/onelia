package routes_transport_http_dto

import (
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
)

type CityDTOResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func CityDTOFromDomain(city domain.City) CityDTOResponse {
	return CityDTOResponse{
		ID:   city.ID,
		Name: city.Name,
	}
}
