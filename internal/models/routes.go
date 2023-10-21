package models

import (
	"github.com/google/uuid"
	"time"
)

type (
	Route struct {
		ID          uuid.UUID `json:"_id" db:"id"`
		Name        string    `json:"name" db:"name"`
		Description string    `json:"description" db:"description"`
		AddressText string    `json:"address_text" db:"address_text"`
		AddressLng  float64   `json:"address_lng" db:"address_lng"`
		AddressLat  float64   `json:"address_lat" db:"address_lat"`
		IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
	}

	ReviewRoute struct {
		ID         uuid.UUID `json:"_id" db:"id"`
		OwnerID    uuid.UUID `json:"owner_id" db:"owner_id"`
		RouteID    uuid.UUID `json:"route_id" db:"route_id"`
		ReviewText string    `json:"review_text" db:"review_text"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		IsDeleted  bool      `json:"is_deleted" db:"is_deleted"`
	}
)
