package routes_transport_http

import (
	"net/http"

	core_logger "github.com/SALutHere/onelia/internal/core/logger"
	core_http_request "github.com/SALutHere/onelia/internal/core/transport/http/request"
	core_http_response "github.com/SALutHere/onelia/internal/core/transport/http/response"
	routes_transport_http_dto "github.com/SALutHere/onelia/internal/routes/transport/http/dto"
)

type GetBestRoutesResponse []routes_transport_http_dto.RouteDTOResponse

// GetBestRoutes godoc
// @Summary		Получить до трёх лучших маршрутов
// @Description	Получить топ 3 лучших маршрутов по времени или по стоимости
// @Tags		routes
// @Produce		json
// @Param		from_city	query	string	true			"Город отправления по маршруту"
// @Param		to_city		query	string	true			"Город прибытия по маршруту"
// @Param		sort		query	string	false			"Критерий определения лучших маршрутов: `time` (по умолчанию) или `price`"
// @Success		200		{object} 	GetBestRoutesResponse				"Успешное получение лучших маршрутов"
// @Failure		400		{object}	core_http_response.ErrorResponse	"Bad request"
// @Failure		404		{object}	core_http_response.ErrorResponse	"City not found"
// @Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
// @Router		/routes		[get]
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
