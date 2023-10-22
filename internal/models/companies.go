package models

import "github.com/google/uuid"

type Company struct {
	ID          uuid.UUID `json:"_id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	IsReleased  bool      `json:"is_released" db:"is_released"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	PhotoCard   string    `json:"photo_card" db:"photo_card"`
}

func (c *Company) IsNil() bool {
	return c.ID.ID() == 0
}
