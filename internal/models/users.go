package models

import (
	"github.com/google/uuid"
	"time"
)

type (
	User struct {
		ID                     uuid.UUID `json:"_id" db:"id"`
		VkID                   int64     `json:"vk_id" db:"vk_id"`
		PassedAppOnboarding    bool      `json:"passed_app_onboarding" db:"passed_app_onboarding"`
		PassedPrismaOnboarding bool      `json:"passed_prisma_onboarding" db:"passed_prisma_onboarding"`
	}

	UserAchievementsRel struct {
		ID            uuid.UUID `json:"_id" db:"id"`
		UserID        uuid.UUID `json:"user_id" db:"user_id"`
		AchievementID uuid.UUID `json:"achievement_id" db:"achievement_id"`
	}

	UserCoinsRel struct {
		ID        uuid.UUID `json:"_id" db:"id"`
		UserID    uuid.UUID `json:"user_id" db:"user_id"`
		Coins     int       `json:"coins" db:"coins"`
		Operation bool      `json:"operation" db:"operation"`
	}

	UserMapFilterRel struct {
		ID       uuid.UUID `json:"_id" db:"id"`
		UserID   uuid.UUID `json:"user_id" db:"user_id"`
		FilterID uuid.UUID `json:"filter_id" db:"filter_id"`
	}

	UserPrivacyRel struct {
		ID                   uuid.UUID `json:"_id" db:"id"`
		UserID               uuid.UUID `json:"user_id" db:"user_id"`
		CanViewAchievements  bool      `json:"can_view_achievements" db:"can_view_achievements"`
		CanViewProgressOnMap bool      `json:"can_view_progress_on_map" db:"can_view_progress_on_map"`
	}

	UserProgressOnMapRel struct {
		ID        uuid.UUID `json:"_id" db:"id"`
		UserID    uuid.UUID `json:"user_id" db:"user_id"`
		RouteID   uuid.UUID `json:"route_id" db:"route_id"`
		EventID   uuid.UUID `json:"event_id" db:"event_id"`
		PlaceID   uuid.UUID `json:"place_id" db:"place_id"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}
)

func (u *User) IsNil() bool {
	return u.ID.ID() == 0
}
