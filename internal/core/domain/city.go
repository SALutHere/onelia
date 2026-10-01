package domain

import "uuid"

type City struct {
	ID   uuid.UUID
	Name string
}

func NewCity(id uuid.UUID, name string) City {
	return City{
		ID:   id,
		Name: name,
	}
}
