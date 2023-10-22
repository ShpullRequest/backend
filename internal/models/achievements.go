package models

import "github.com/google/uuid"

type (
	Achievements struct {
		ID          uuid.UUID `json:"_id" db:"id"`
		Name        string    `json:"name" db:"name"`
		Description string    `json:"description" db:"description"`
		Icon        string    `json:"icon" db:"icon"`
		Coins       int       `json:"coins" db:"coins"`
	}
)
