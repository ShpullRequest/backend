package repository

import (
	"context"

	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewPlace(ctx context.Context, place models.Place) (*models.Place, error) {
	id, err := p.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO places (name, description, carousel, address_text, address_lng, address_lat, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		place.Name,
		place.Description,
		place.Carousel,
		place.AddressText,
		place.AddressLng,
		place.AddressLat,
		place.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	place.ID = id
	return &place, nil
}

func (p *Pg) GetPlace(ctx context.Context, id uuid.UUID) (*models.Place, error) {
	var place models.Place
	err := p.db.GetContext(ctx, &place, "SELECT * FROM places WHERE id = $1", id)

	return &place, err
}

func (p *Pg) NewReviewPlace(ctx context.Context, reviewPlace models.ReviewPlace) (*models.ReviewPlace, error) {
	id, err := p.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO reviews_places (owner_id, place_id, review_text, stars, created_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6)",
		reviewPlace.OwnerID,
		reviewPlace.PlaceID,
		reviewPlace.ReviewText,
		reviewPlace.Stars,
		reviewPlace.CreatedAt,
		reviewPlace.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	reviewPlace.ID = id
	return &reviewPlace, nil
}

func (p *Pg) GetReviewsPlace(ctx context.Context, placeID uuid.UUID) ([]models.ReviewPlace, error) {
	var reviewsPlace []models.ReviewPlace
	err := p.db.SelectContext(ctx, &reviewsPlace, "SELECT * FROM reviews_places WHERE place_id = $1", placeID)

	return reviewsPlace, err
}
