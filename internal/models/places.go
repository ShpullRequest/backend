package models

import (
	"github.com/google/uuid"
	"time"
)

type (
	Place struct {
		ID          uuid.UUID `json:"_id" db:"id"`
		Name        string    `json:"name" db:"name"`
		Description string    `json:"description" db:"description"`
		Carousel    []string  `json:"carousel" db:"carousel"`
		AddressText string    `json:"address_text" db:"address_text"`
		AddressLng  float64   `json:"address_lng" db:"address_lng"`
		AddressLat  float64   `json:"address_lat" db:"address_lat"`
		IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
	}

	ReviewPlace struct {
		ID         uuid.UUID `json:"_id" db:"id"`
		OwnerID    uuid.UUID `json:"owner_id" db:"owner_id"`
		PlaceID    uuid.UUID `json:"place_id" db:"place_id"`
		ReviewText string    `json:"review_text" db:"review_text"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		IsDeleted  bool      `json:"is_deleted" db:"is_deleted"`
	}
)
