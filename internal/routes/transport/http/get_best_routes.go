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

	fromCity, err := core_http_request.GetRequiredStringQueryParam(r, "from_city")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get origin city query parameter",
		)
		return
	}

	toCity, err := core_http_request.GetRequiredStringQueryParam(r, "to_city")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get destination city query parameter",
		)
		return
	}

	sort := core_http_request.GetSortTypeQueryParam(r, "sort")

	routes, err := h.routesService.GetBestRoutes(ctx, fromCity, toCity, sort)
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
