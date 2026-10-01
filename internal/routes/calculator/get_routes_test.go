package routes_calculator

import (
	"context"
	"errors"
	"testing"
	"uuid"

	"github.com/SALutHere/onelia/internal/core/domain"
)

// Cancel on the first recursive step to exercise propagation from inside walk.
type cancelDuringWalkContext struct {
	context.Context
	cancel context.CancelFunc
	calls  int
}

func (c *cancelDuringWalkContext) Err() error {
	c.calls++
	if c.calls == 3 {
		c.cancel()
	}
	return c.Context.Err()
}

func TestGetRoutesPropagatesCancellationDuringWalk(t *testing.T) {
	from := domain.NewCity(uuid.New(), "Moscow")
	to := domain.NewCity(uuid.New(), "Saint Petersburg")
	calculator := &RoutesCalculator{
		cities:       map[uuid.UUID]domain.City{from.ID: from, to.ID: to},
		citiesByName: map[string]uuid.UUID{from.Name: from.ID, to.Name: to.ID},
		outgoing: map[uuid.UUID][]domain.Segment{
			from.ID: {domain.NewSegment(uuid.New(), from.ID, to.ID, "train", 240, 450000)},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	routes, err := calculator.GetRoutes(&cancelDuringWalkContext{Context: ctx, cancel: cancel}, from.Name, to.Name, 3)
	if !errors.Is(err, context.Canceled) || routes != nil {
		t.Fatalf("routes = %v, error = %v; want nil routes and context.Canceled", routes, err)
	}
}
