package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type (
	Route struct {
		ID          uuid.UUID      `json:"_id" db:"id"`
		CompanyID   *uuid.UUID     `json:"company_id,omitempty" db:"company_id"`
		Name        string         `json:"name" db:"name"`
		Description string         `json:"description" db:"description"`
		Places      pq.StringArray `json:"-" db:"places"`
		Events      pq.StringArray `json:"-" db:"events"`
		IsDeleted   bool           `json:"-" db:"is_deleted"`
	}

	RouteGeo struct {
		Type   string      `json:"type"`
		Object interface{} `json:"object"`
	}

	RouteWithGeo struct {
		Route
		Geo []RouteGeo `json:"geo"`
	}

	ReviewRoute struct {
		ID         uuid.UUID `json:"_id" db:"id"`
		OwnerID    uuid.UUID `json:"owner_id" db:"owner_id"`
		RouteID    uuid.UUID `json:"route_id" db:"route_id"`
		ReviewText string    `json:"review_text" db:"review_text"`
		Stars      float64   `json:"stars" db:"stars"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		IsDeleted  bool      `json:"is_deleted" db:"is_deleted"`
	}
)
