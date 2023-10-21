package repository

import (
	"context"

	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewPlace(ctx context.Context, place models.Place) (*models.Place, error) {
	_, err := p.db.ExecContext(
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

	return &place, nil
}

func (p *Pg) GetPlace(ctx context.Context, id uuid.UUID) (*models.Place, error) {
	var place models.Place
	err := p.db.GetContext(ctx, &place, "SELECT * FROM places WHERE id = $1", id)

	return &place, err
}

func (p *Pg) NewReviewPlace(ctx context.Context, reviewPlace models.ReviewPlace) (*models.ReviewPlace, error) {
	_, err := p.db.ExecContext(
		ctx,
		"INSERT INTO reviews_places (owner_id, place_id, review_text, created_at, is_deleted) VALUES ($1, $2, $3, $4, $5)",
		reviewPlace.OwnerID,
		reviewPlace.PlaceID,
		reviewPlace.ReviewText,
		reviewPlace.CreatedAt,
		reviewPlace.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &reviewPlace, nil
}

func (p *Pg) GetReviewsPlace(ctx context.Context, placeID uuid.UUID) ([]models.ReviewPlace, error) {
	var reviewsPlace []models.ReviewPlace
	err := p.db.SelectContext(ctx, &reviewsPlace, "SELECT * FROM reviews_places WHERE place_id = $1", placeID)

	return reviewsPlace, err
}
