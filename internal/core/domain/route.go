package domain

type Route struct {
	Parts                []RoutePart
	TotalDurationMinutes int
	TotalPrice           int64
}

type RoutePart struct {
	Segment
	FromCity City
	ToCity   City
}
