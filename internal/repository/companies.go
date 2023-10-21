package repository

import (
	"context"

	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewCompany(ctx context.Context, company models.Company) (*models.Company, error) {
	_, err := p.db.ExecContext(
		ctx,
		"INSERT INTO companies (user_id, is_organisation, name, description, photo_card) VALUES ($1, $2, $3, $4, $5)",
		company.UserID,
		company.IsOrganisation,
		company.Name,
		company.Description,
		company.PhotoCard,
	)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (p *Pg) GetCompanyByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	var company models.Company
	err := p.db.GetContext(ctx, &company, "SELECT * FROM companies WHERE id = $1", id)

	return &company, err
}
