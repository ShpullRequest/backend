package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type (
	Place struct {
		ID          uuid.UUID      `json:"_id" db:"id"`
		Name        string         `json:"name" db:"name"`
		Description string         `json:"description" db:"description"`
		Carousel    pq.StringArray `json:"carousel" db:"carousel" swaggertype:"array,string"`
		AddressText string         `json:"address_text" db:"address_text"`
		AddressLng  float64        `json:"address_lng" db:"address_lng"`
		AddressLat  float64        `json:"address_lat" db:"address_lat"`
		IsDeleted   bool           `json:"is_deleted" db:"is_deleted"`
	}

	ReviewPlace struct {
		ID         uuid.UUID `json:"_id" db:"id"`
		OwnerID    uuid.UUID `json:"owner_id" db:"owner_id"`
		PlaceID    uuid.UUID `json:"place_id" db:"place_id"`
		ReviewText string    `json:"review_text" db:"review_text"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		Stars      float64   `json:"stars" db:"stars"`
		IsDeleted  bool      `json:"is_deleted" db:"is_deleted"`
	}
)
