package repository

import (
	"context"
	"fmt"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewRoute(ctx context.Context, route models.Route) (*models.Route, error) {
	id, err := p.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO routes (company_id, name, description, places, events, is_deleted) VALUES ($1, $2, $3, $4, $5, $6)",
		route.CompanyID,
		route.Name,
		route.Description,
		route.Places,
		route.Events,
		route.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	route.ID = id
	return &route, nil
}

func (p *Pg) GetRoute(ctx context.Context, id uuid.UUID) (*models.RouteWithGeo, error) {
	var route models.Route
	err := p.db.GetContext(ctx, &route, "SELECT * FROM routes WHERE id = $1", id)

	routeWithGeo := p.routeToRouteWithGeo(ctx, route)
	return &routeWithGeo, err
}

func (p *Pg) SearchRoutes(ctx context.Context, q string) ([]models.RouteWithGeo, error) {
	q = fmt.Sprintf("%%%s%%", q)

	var routes []models.Route
	err := p.db.SelectContext(
		ctx,
		&routes,
		`SELECT * FROM routes
				WHERE 
					LOWER(name) LIKE LOWER($1) OR
					LOWER(description) LIKE LOWER($1)
				`,
		q,
	)

	return p.sliceRouteToSliceRouteWithGeo(ctx, routes), err
}

func (p *Pg) GetAllRoutesByCompanyID(ctx context.Context, companyID uuid.UUID) ([]models.RouteWithGeo, error) {
	var routes []models.Route
	err := p.db.SelectContext(ctx, &routes, "SELECT * FROM routes WHERE company_id = $1", companyID)

	return p.sliceRouteToSliceRouteWithGeo(ctx, routes), err
}

func (p *Pg) GetAllRoutes(ctx context.Context) ([]models.RouteWithGeo, error) {
	var routes []models.Route
	err := p.db.SelectContext(ctx, &routes, "SELECT * FROM routes WHERE is_deleted = false")

	return p.sliceRouteToSliceRouteWithGeo(ctx, routes), err
}

func (p *Pg) SaveRoute(ctx context.Context, routeWithGeo *models.RouteWithGeo) error {
	route := routeWithGeo.Route

	_, err := p.db.ExecContext(
		ctx,
		`
			UPDATE routes 
				SET name = $1, description = $2, events = $3, places = $4 
				WHERE id = $5
		`,
		route.Name, route.Description, route.Events, route.Places,
		route.ID,
	)

	routeWithGeo.Geo = p.routeToRouteWithGeo(ctx, route).Geo

	return err
}

func (p *Pg) NewReviewRoute(ctx context.Context, reviewRoute models.ReviewRoute) (*models.ReviewRoute, error) {
	id, err := p.db.ExecContextWithReturnID(
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

	reviewRoute.ID = id
	return &reviewRoute, nil
}

func (p *Pg) SaveReviewRoute(ctx context.Context, reviewEvent *models.ReviewRoute) error {
	_, err := p.db.ExecContext(
		ctx,
		`UPDATE reviews_routes SET review_text = $1, stars = $2 WHERE owner_id = $3 AND route_id = $4`,
		reviewEvent.ReviewText, reviewEvent.Stars, reviewEvent.OwnerID, reviewEvent.RouteID,
	)

	return err
}

func (p *Pg) GetReviewRoute(ctx context.Context, ownerID uuid.UUID, routeID uuid.UUID) (*models.ReviewRoute, error) {
	var reviewRoute models.ReviewRoute
	err := p.db.GetContext(ctx, &reviewRoute, "SELECT * FROM reviews_routes WHERE owner_id = $1 AND route_id = $2 AND is_deleted = false", ownerID, routeID)

	return &reviewRoute, err
}

func (p *Pg) GetReviewsRoute(ctx context.Context, routeID uuid.UUID) ([]models.ReviewRoute, error) {
	var reviewRoute []models.ReviewRoute
	err := p.db.SelectContext(ctx, &reviewRoute, "SELECT * FROM reviews_routes WHERE route_id = $1 AND is_deleted = false", routeID)

	return reviewRoute, err
}

func (p *Pg) sliceRouteToSliceRouteWithGeo(ctx context.Context, routes []models.Route) []models.RouteWithGeo {
	var routesWithGeo []models.RouteWithGeo

	for _, route := range routes {
		routesWithGeo = append(routesWithGeo, p.routeToRouteWithGeo(ctx, route))
	}

	return routesWithGeo
}

func (p *Pg) routeToRouteWithGeo(ctx context.Context, route models.Route) models.RouteWithGeo {
	rWithGeo := models.RouteWithGeo{Route: route}

	for _, e := range route.Events {
		eID, err := uuid.Parse(e)
		if err != nil {
			continue
		}

		event, err := p.GetEvent(ctx, eID)
		if err != nil {
			continue
		}

		rWithGeo.Geo = append(rWithGeo.Geo, models.RouteGeo{
			Type:   "event",
			Object: event,
		})
	}

	for _, pl := range route.Places {
		pID, err := uuid.Parse(pl)
		if err != nil {
			continue
		}

		place, err := p.GetPlace(ctx, pID)
		if err != nil {
			continue
		}

		rWithGeo.Geo = append(rWithGeo.Geo, models.RouteGeo{
			Type:   "place",
			Object: place,
		})
	}

	return rWithGeo
}
