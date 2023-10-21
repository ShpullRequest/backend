package repository

import (
	"context"

	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewRoute(ctx context.Context, route models.Route) (*models.Route, error) {
	id, err := p.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO routes (company_id, name, description, address_text, address_lng, address_lat, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		route.Name,
		route.CompanyID,
		route.Description,
		route.AddressText,
		route.AddressLng,
		route.AddressLat,
		route.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	route.ID = id
	return &route, nil
}

func (p *Pg) GetRoute(ctx context.Context, id uuid.UUID) (*models.Route, error) {
	var route models.Route
	err := p.db.GetContext(ctx, &route, "SELECT * FROM routes WHERE id = $1", id)

	return &route, err
}

func (p *Pg) NewReviewRoute(ctx context.Context, reviewRoute models.ReviewRoute) (*models.ReviewRoute, error) {
	_, err := p.db.ExecContext(
		ctx,
		"INSERT INTO reviews_routes (owner_id, route_id, review_text, stars, created_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6)",
		reviewRoute.OwnerID,
		reviewRoute.RouteID,
		reviewRoute.ReviewText,
		reviewRoute.Stars,
		reviewRoute.CreatedAt,
		reviewRoute.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &reviewRoute, nil
}

func (p *Pg) GetReviewsRoute(ctx context.Context, routeID uuid.UUID) ([]models.ReviewRoute, error) {
	var reviewRoute []models.ReviewRoute
	err := p.db.SelectContext(ctx, &reviewRoute, "SELECT * FROM reviews_routes WHERE route_id = $1", routeID)

	return reviewRoute, err
}
