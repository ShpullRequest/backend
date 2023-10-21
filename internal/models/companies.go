package models

import "github.com/google/uuid"

type Company struct {
	ID             uuid.UUID `json:"_id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	IsOrganisation bool      `json:"is_organisation" db:"is_organisation"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	PhotoCard      string    `json:"photo_card" db:"photo_card"`
}
