package repository

import (
	"context"

	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewAchievement(ctx context.Context, achievement models.Achievements) (*models.Achievements, error) {
	_, err := p.db.ExecContext(
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

	return &achievement, nil
}

func (p *Pg) GetAchievementByID(ctx context.Context, id uuid.UUID) (*models.Achievements, error) {
	var achievement models.Achievements
	err := p.db.GetContext(ctx, &achievement, "SELECT * FROM achievements WHERE id = $1", id)

	return &achievement, err
}
