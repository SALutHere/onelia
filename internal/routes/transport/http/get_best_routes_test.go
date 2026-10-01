package routes_transport_http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
	core_logger "github.com/SALutHere/onelia/internal/core/logger"
	core_http_server "github.com/SALutHere/onelia/internal/core/transport/http/server"
	routes_calculator "github.com/SALutHere/onelia/internal/routes/calculator"
	routes_service "github.com/SALutHere/onelia/internal/routes/service"
	routes_transport_http "github.com/SALutHere/onelia/internal/routes/transport/http"
	routes_transport_http_dto "github.com/SALutHere/onelia/internal/routes/transport/http/dto"
	"go.uber.org/zap"
)

type testRepository struct {
	cities   []domain.City
	segments []domain.Segment
}

func (r testRepository) GetCities(context.Context) ([]domain.City, error) {
	return r.cities, nil
}

func (r testRepository) GetSegments(context.Context) ([]domain.Segment, error) {
	return r.segments, nil
}

func TestGetBestRoutesByCityName(t *testing.T) {
	names := []string{"Moscow", "Saint Petersburg", "Kazan", "Ufa", "Perm", "Санкт-Петербург"}
	repository := testRepository{}
	for _, name := range names {
		repository.cities = append(repository.cities, domain.NewCity(uuid.New(), name))
	}
	addSegment := func(from, to, minutes int, price int64) {
		repository.segments = append(repository.segments, domain.NewSegment(
			uuid.New(), repository.cities[from].ID, repository.cities[to].ID,
			"train", minutes, price,
		))
	}
	addSegment(0, 1, 50, 500)
	addSegment(0, 1, 100, 100)
	addSegment(0, 1, 150, 50)
	addSegment(0, 2, 20, 200)
	addSegment(2, 1, 40, 100)
	addSegment(2, 3, 10, 10)
	addSegment(3, 1, 10, 10)
	// This faster, cheaper path requires three transfers and must be excluded.
	addSegment(3, 4, 1, 1)
	addSegment(4, 1, 1, 1)
	// Revisiting the origin must also be excluded.
	addSegment(2, 0, 1, 1)

	calculator, err := routes_calculator.NewRoutesCalculator(context.Background(), repository)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("all routes obey transfer and cycle limits", func(t *testing.T) {
		routes, err := calculator.GetRoutes(context.Background(), "Moscow", "Saint Petersburg", 3)
		if err != nil {
			t.Fatal(err)
		}
		if len(routes) != 5 {
			t.Fatalf("got %d routes, want 5", len(routes))
		}
		for _, route := range routes {
			if len(route.Parts) > 3 {
				t.Fatalf("too many segments: %+v", route)
			}
			visited := map[uuid.UUID]bool{repository.cities[0].ID: true}
			for _, part := range route.Parts {
				if visited[part.ToCityID] {
					t.Fatalf("repeated city: %+v", route)
				}
				visited[part.ToCityID] = true
			}
		}
	})
	handler := routes_transport_http.NewRoutesHTTPHandler(routes_service.NewRoutesService(calculator))
	router := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	router.RegisterRoutes(handler.Routes()...)
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", router))
	logger := &core_logger.Logger{Logger: zap.NewNop()}

	tests := []struct {
		name       string
		query      url.Values
		status     int
		durations  []int
		prices     []int64
		errorValue string
	}{
		{name: "default time order", query: url.Values{"from_city": {"Moscow"}, "to_city": {"Saint Petersburg"}}, status: 200, durations: []int{40, 50, 60}, prices: []int64{220, 500, 300}},
		{name: "explicit time order", query: url.Values{"from_city": {"Moscow"}, "to_city": {"Saint Petersburg"}, "sort": {"time"}}, status: 200, durations: []int{40, 50, 60}, prices: []int64{220, 500, 300}},
		{name: "price order", query: url.Values{"from_city": {"Moscow"}, "to_city": {"Saint Petersburg"}, "sort": {"price"}}, status: 200, durations: []int{150, 100, 40}, prices: []int64{50, 100, 220}},
		{name: "surrounding whitespace", query: url.Values{"from_city": {"  Moscow\t"}, "to_city": {" Saint Petersburg "}}, status: 200, durations: []int{40, 50, 60}, prices: []int64{220, 500, 300}},
		{name: "missing origin", query: url.Values{"to_city": {"Saint Petersburg"}}, status: 400, errorValue: "from_city"},
		{name: "missing destination", query: url.Values{"from_city": {"Moscow"}}, status: 400, errorValue: "to_city"},
		{name: "blank origin", query: url.Values{"from_city": {" \t"}, "to_city": {"Saint Petersburg"}}, status: 400, errorValue: "from_city"},
		{name: "blank destination", query: url.Values{"from_city": {"Moscow"}, "to_city": {" "}}, status: 400, errorValue: "to_city"},
		{name: "same city", query: url.Values{"from_city": {"Moscow"}, "to_city": {" Moscow "}}, status: 400, errorValue: "must differ"},
		{name: "unknown origin", query: url.Values{"from_city": {"Unknown"}, "to_city": {"Saint Petersburg"}}, status: 404, errorValue: "Unknown"},
		{name: "unknown destination", query: url.Values{"from_city": {"Moscow"}, "to_city": {"Unknown"}}, status: 404, errorValue: "Unknown"},
		{name: "case sensitive names", query: url.Values{"from_city": {"moscow"}, "to_city": {"Saint Petersburg"}}, status: 404, errorValue: "moscow"},
		{name: "invalid sort", query: url.Values{"from_city": {"Moscow"}, "to_city": {"Saint Petersburg"}, "sort": {"invalid"}}, status: 400, errorValue: "sort"},
		{name: "UUID is not a city name", query: url.Values{"from_city": {repository.cities[0].ID.String()}, "to_city": {"Saint Petersburg"}}, status: 404, errorValue: repository.cities[0].ID.String()},
		{name: "old query parameters", query: url.Values{"from_city_id": {repository.cities[0].ID.String()}, "to_city_id": {repository.cities[1].ID.String()}}, status: 400, errorValue: "from_city"},
		{name: "no routes", query: url.Values{"from_city": {"Saint Petersburg"}, "to_city": {"Moscow"}}, status: 200},
		{name: "Unicode city name", query: url.Values{"from_city": {"Moscow"}, "to_city": {"Санкт-Петербург"}}, status: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/routes?"+tt.query.Encode(), nil)
			req = req.WithContext(core_logger.ToContext(req.Context(), logger))
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, req)
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, tt.status, response.Body.String())
			}
			if tt.status != http.StatusOK {
				var body struct{ Error string }
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(body.Error, tt.errorValue) {
					t.Fatalf("error = %q, want it to contain %q", body.Error, tt.errorValue)
				}
				return
			}
			var routes []routes_transport_http_dto.RouteDTOResponse
			if err := json.Unmarshal(response.Body.Bytes(), &routes); err != nil {
				t.Fatal(err)
			}
			if routes == nil || len(routes) != len(tt.durations) {
				t.Fatalf("want array of %d routes; body: %s", len(tt.durations), response.Body.String())
			}
			for i, route := range routes {
				if route.TotalDurationMinutes != tt.durations[i] || route.TotalPrice != tt.prices[i] {
					t.Fatalf("route %d totals = (%d, %d), want (%d, %d)", i, route.TotalDurationMinutes, route.TotalPrice, tt.durations[i], tt.prices[i])
				}
				if len(route.Parts) < 1 || len(route.Parts) > 3 {
					t.Fatalf("route has %d segments", len(route.Parts))
				}
				previous := repository.cities[0]
				visited := map[uuid.UUID]bool{previous.ID: true}
				for _, part := range route.Parts {
					if part.FromCity.ID != previous.ID || part.FromCity.Name != previous.Name || visited[part.ToCity.ID] {
						t.Fatalf("disconnected route or repeated city: %+v", route)
					}
					visited[part.ToCity.ID] = true
					previous = domain.NewCity(part.ToCity.ID, part.ToCity.Name)
				}
				if previous != repository.cities[1] {
					t.Fatalf("wrong destination: %+v", previous)
				}
			}
		})
	}
}
