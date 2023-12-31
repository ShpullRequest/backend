package repository

import (
	"context"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewUser(ctx context.Context, user models.User) (*models.User, error) {
	id, err := p.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO users (vk_id, passed_onboarding, selected_geo) VALUES ($1, $2, $3)",
		user.VkID,
		user.PassedOnboarding,
		user.SelectedGeo,
	)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return &user, err
}

func (p *Pg) SaveUser(ctx context.Context, user *models.User) error {
	_, err := p.db.ExecContext(
		ctx,
		"UPDATE users SET passed_onboarding = $1, selected_geo = $2 WHERE id = $3",
		user.PassedOnboarding,
		user.SelectedGeo,
		user.ID,
	)

	return err
}

func (p *Pg) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := p.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)

	return &user, err
}

func (p *Pg) GetUserByVkID(ctx context.Context, vkID int64) (*models.User, error) {
	var user models.User
	err := p.db.GetContext(ctx, &user, "SELECT * FROM users WHERE vk_id = $1", vkID)

	return &user, err
}
