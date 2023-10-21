package models

import (
	"github.com/google/uuid"
	"time"
)

type (
	Event struct {
		ID          uuid.UUID `json:"_id" db:"id"`
		CompanyID   uuid.UUID `json:"company_id" db:"company_id"`
		Name        string    `json:"name" db:"name"`
		Description string    `json:"description" db:"description"`
		Carousel    []string  `json:"carousel" db:"carousel"`
		Icon        string    `json:"icon" db:"icon"`
		StartTime   time.Time `json:"start_time" db:"start_time"`
		AddressText string    `json:"address_text" db:"address_text"`
		AddressLng  float64   `json:"address_lng" db:"address_lng"`
		AddressLat  float64   `json:"address_lat" db:"address_lat"`
		IsDeleted   bool      `json:"-" db:"is_deleted"`
	}

	EventsFilter struct {
		ID   uuid.UUID `json:"_id" db:"id"`
		Name string    `json:"name" db:"name"`
	}

	EventsFiltersRel struct {
		ID       uuid.UUID `json:"_id" db:"id"`
		EventID  uuid.UUID `json:"event_id" db:"event_id"`
		FilterID uuid.UUID `json:"filter_id" db:"filter_id"`
	}

	ReviewEvent struct {
		ID         uuid.UUID `json:"_id" db:"id"`
		OwnerID    uuid.UUID `json:"owner_id" db:"owner_id"`
		EventID    uuid.UUID `json:"event_id" db:"event_id"`
		ReviewText string    `json:"review_text" db:"review_text"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		IsDeleted  bool      `json:"is_deleted" db:"is_deleted"`
	}
)
