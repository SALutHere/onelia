package routes_transport_http

import (
	"net/http"

	core_logger "github.com/SALutHere/onelia/internal/core/logger"
	core_http_request "github.com/SALutHere/onelia/internal/core/transport/http/request"
	core_http_response "github.com/SALutHere/onelia/internal/core/transport/http/response"
	routes_transport_http_dto "github.com/SALutHere/onelia/internal/routes/transport/http/dto"
)

type GetBestRoutesResponse []routes_transport_http_dto.RouteDTOResponse

func (h *RoutesHTTPHandler) GetBestRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	fromCityID, err := core_http_request.GetRequiredUUIDQueryParam(r, "from_city_id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get origin city path value",
		)
		return
	}

	toCityID, err := core_http_request.GetRequiredUUIDQueryParam(r, "to_city_id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get destination city path value",
		)
		return
	}

	sort := core_http_request.GetSortTypeQueryParam(r, "sort")

	routes, err := h.routesService.GetBestRoutes(ctx, fromCityID, toCityID, sort)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get routes",
		)
		return
	}

	response := GetBestRoutesResponse(routes_transport_http_dto.RoutesDTOFromDomains(routes))

	responseHandler.JSONResponse(response, http.StatusOK)
}
