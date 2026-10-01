package routes_transport_http

import (
	"net/http"

	core_http_server "github.com/SALutHere/onelia/internal/core/transport/http/server"
)

func (h *RoutesHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/routes",
			Handler: h.GetBestRoutes,
		},
	}
}
