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

func CitiesDTOFromDomains(cities []domain.City) []CityDTOResponse {
	citiesDTO := make([]CityDTOResponse, len(cities))

	for i, city := range cities {
		citiesDTO[i] = CityDTOFromDomain(city)
	}

	return citiesDTO
}
