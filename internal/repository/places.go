package repository

import (
	"context"
	"fmt"

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

func (p *Pg) SearchPlace(ctx context.Context, q string) ([]models.Place, error) {
	q = fmt.Sprintf("%%%s%%", q)

	var places []models.Place
	err := p.db.SelectContext(
		ctx,
		&places,
		`SELECT * FROM places
				WHERE 
					LOWER(name) LIKE LOWER($1) OR
					LOWER(description) LIKE LOWER($1) OR
					LOWER(address_text) LIKE LOWER($1);
				`,
		q,
	)

	return places, err
}

func (p *Pg) GetAllPlaces(ctx context.Context) ([]models.Place, error) {
	var places []models.Place
	err := p.db.SelectContext(ctx, &places, "SELECT * FROM places")

	return places, err
}

func (p *Pg) SavePlace(ctx context.Context, place *models.Place) error {
	_, err := p.db.ExecContext(
		ctx,
		"UPDATE places SET name = $1, description = $2, carousel = $3, address_text = $4, address_lng = $5, address_lat = $6, is_deleted = $7 WHERE id = $8",
		place.Name,
		place.Description,
		place.Carousel,
		place.AddressText,
		place.AddressLng,
		place.AddressLat,
		place.IsDeleted,
		place.ID,
	)

	return err
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

func (p *Pg) SaveReviewPlace(ctx context.Context, reviewPlace *models.ReviewEvent) error {
	_, err := p.db.ExecContext(
		ctx,
		`UPDATE reviews_places SET review_text = $1, stars = $2 WHERE owner_id = $3 AND place_id = $4`,
		reviewPlace.ReviewText, reviewPlace.Stars, reviewPlace.OwnerID, reviewPlace.EventID,
	)

	return err
}

func (p *Pg) GetReviewPlace(ctx context.Context, ownerID uuid.UUID, placeID uuid.UUID) (*models.ReviewEvent, error) {
	var reviewEvent models.ReviewEvent
	err := p.db.GetContext(ctx, &reviewEvent, "SELECT * FROM reviews_places WHERE owner_id = $1 AND place_id = $2 AND is_deleted = false", ownerID, placeID)

	return &reviewEvent, err
}

func (p *Pg) GetReviewsPlace(ctx context.Context, placeID uuid.UUID) ([]models.ReviewPlace, error) {
	var reviewsPlace []models.ReviewPlace
	err := p.db.SelectContext(ctx, &reviewsPlace, "SELECT * FROM reviews_places WHERE place_id = $1", placeID)

	return reviewsPlace, err
}
