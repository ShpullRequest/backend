package repository

import (
	"context"

	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewAchievement(ctx context.Context, achievement models.Achievements) (*models.Achievements, error) {
	id, err := p.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO achievements (name, description, icon, coins) VALUES ($1, $2, $3, $4)",
		achievement.Name,
		achievement.Description,
		achievement.Icon,
		achievement.Coins,
	)
	if err != nil {
		return nil, err
	}

	achievement.ID = id
	return &achievement, nil
}

func (p *Pg) GetAchievementByID(ctx context.Context, id uuid.UUID) (*models.Achievements, error) {
	var achievement models.Achievements
	err := p.db.GetContext(ctx, &achievement, "SELECT * FROM achievements WHERE id = $1", id)

	return &achievement, err
}

func (p *Pg) GetAllAchievements(ctx context.Context) ([]models.Achievements, error) {
	var achievements []models.Achievements
	err := p.db.SelectContext(ctx, &achievements, "SELECT * FROM achievements")

	return achievements, err
}

func (p *Pg) SaveAchievement(ctx context.Context, achievement *models.Achievements) error {
	_, err := p.db.ExecContext(
		ctx,
		"UPDATE achievements SET name = $1, description = $2, icon = $3, coins = $4 WHERE id = $5",
		achievement.Name,
		achievement.Description,
		achievement.Icon,
		achievement.Coins,
		achievement.ID,
	)

	return err
}
