package repository

import (
	"context"

	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewCompany(ctx context.Context, company models.Company) (*models.Company, error) {
	id, err := p.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO companies (user_id, name, description, photo_card) VALUES ($1, $2, $3, $4)",
		company.UserID,
		company.Name,
		company.Description,
		company.PhotoCard,
	)
	if err != nil {
		return nil, err
	}

	company.ID = id
	return &company, err
}

func (p *Pg) GetCompanyByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	var company models.Company
	err := p.db.GetContext(ctx, &company, "SELECT * FROM companies WHERE id = $1", id)

	return &company, err
}

func (p *Pg) GetCompanyAverageRating(ctx context.Context, id uuid.UUID) (float64, error) {
	var averageRating float64
	err := p.db.GetContext(
		ctx,
		&averageRating,
		`SELECT calculate_company_rating($1) AS rating`,
		id,
	)

	return averageRating, err
}

type companyWithRating struct {
	models.Company
	Rating float64 `json:"rating" bd:"rating"`
}

func (p *Pg) GetAllCompanies(ctx context.Context) ([]companyWithRating, error) {
	var companies []companyWithRating
	err := p.db.SelectContext(
		ctx,
		&companies,
		`SELECT *, calculate_company_rating(c.id) AS rating FROM companies c WHERE is_released = true`,
	)

	return companies, err
}

func (p *Pg) GetCompaniesByVkID(ctx context.Context, vkID int64) ([]companyWithRating, error) {
	var companies []companyWithRating
	err := p.db.SelectContext(ctx, &companies, "SELECT *, calculate_company_rating(c.id) AS rating FROM companies c WHERE user_id = (SELECT id FROM users WHERE vk_id = $1)", vkID)

	return companies, err
}

func (p *Pg) SaveCompany(ctx context.Context, company *models.Company) error {
	_, err := p.db.ExecContext(
		ctx,
		"UPDATE companies SET is_released = $1, name = $2, description = $3, photo_card = $4  WHERE id = $5",
		company.IsReleased,
		company.Name,
		company.Description,
		company.PhotoCard,
		company.ID,
	)

	return err
}
