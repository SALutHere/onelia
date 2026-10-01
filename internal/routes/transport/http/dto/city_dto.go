package routes_transport_http_dto

import (
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
)

type CityDTOResponse struct {
	ID   uuid.UUID `json:"id"   example:"fa1c8a32-2229-4ba7-843a-0196ebc1caa7"`
	Name string    `json:"name" example:"Berlin"`
}

func CityDTOFromDomain(city domain.City) CityDTOResponse {
	return CityDTOResponse{
		ID:   city.ID,
		Name: city.Name,
	}
}
