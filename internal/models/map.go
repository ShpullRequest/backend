package models

import "github.com/google/uuid"

type MapFilter struct {
	ID   uuid.UUID `json:"_id" db:"id"`
	Name string    `json:"name" db:"name"`
}
