package repository

import (
	"context"
	"fmt"

	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewEvent(ctx context.Context, event models.Event) (*models.Event, error) {
	id, err := p.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO events (company_id, name, description, carousel, tags, icon, start_time, address_text, address_lng, address_lat, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		event.CompanyID,
		event.Name,
		event.Description,
		event.Carousel,
		event.Tags,
		event.Icon,
		event.StartTime,
		event.AddressText,
		event.AddressLng,
		event.AddressLat,
		event.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	event.ID = id
	return &event, nil
}

func (p *Pg) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	var event models.Event
	err := p.db.GetContext(ctx, &event, "SELECT * FROM events WHERE id = $1", id)

	return &event, err
}

func (p *Pg) SearchEvents(ctx context.Context, q string) ([]models.Event, error) {
	q = fmt.Sprintf("%%%s%%", q)

	var events []models.Event
	err := p.db.SelectContext(
		ctx,
		&events,
		`SELECT * FROM events
				WHERE 
					LOWER(name) LIKE LOWER($1) OR
					LOWER(description) LIKE LOWER($1) OR
					LOWER(address_text) LIKE LOWER($1);
				`,
		q,
	)

	return events, err
}

func (p *Pg) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	err := p.db.SelectContext(ctx, &events, "SELECT * FROM events")

	return events, err
}

func (p *Pg) GetEventByCompanyID(ctx context.Context, eventID uuid.UUID, companyID uuid.UUID) (*models.Event, error) {
	var event models.Event
	err := p.db.GetContext(ctx, &event, "SELECT * FROM events WHERE id = $1 AND company_id = $2", eventID, companyID)

	return &event, err
}

func (p *Pg) GetAllEventsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]models.Event, error) {
	var events []models.Event
	err := p.db.SelectContext(ctx, &events, "SELECT * FROM events WHERE company_id = $1", companyID)

	return events, err
}

func (p *Pg) SaveEvent(ctx context.Context, event *models.Event) error {
	_, err := p.db.ExecContext(
		ctx,
		`
			UPDATE events 
				SET name = $1, description = $2, carousel = $3, tags = $4,
				    icon = $5, start_time = $6, address_text = $7, 
				    address_lng = $8, address_lat = $9, is_deleted = $10  
				WHERE id = $11
		`,
		event.Name, event.Description, event.Carousel, event.Tags,
		event.Icon, event.StartTime, event.AddressText,
		event.AddressLng, event.AddressLat, event.IsDeleted,
		event.ID,
	)

	return err
}

func (p *Pg) NewReviewEvent(ctx context.Context, reviewEvent models.ReviewEvent) (*models.ReviewEvent, error) {
	id, err := p.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO reviews_events (owner_id, event_id, review_text, stars, created_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6)",
		reviewEvent.OwnerID,
		reviewEvent.EventID,
		reviewEvent.ReviewText,
		reviewEvent.Stars,
		reviewEvent.CreatedAt,
		reviewEvent.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	reviewEvent.ID = id
	return &reviewEvent, nil
}

func (p *Pg) SaveReviewEvent(ctx context.Context, reviewEvent *models.ReviewEvent) error {
	_, err := p.db.ExecContext(
		ctx,
		`UPDATE reviews_events SET review_text = $1, stars = $2 WHERE owner_id = $3 AND event_id = $4`,
		reviewEvent.ReviewText, reviewEvent.Stars, reviewEvent.OwnerID, reviewEvent.EventID,
	)

	return err
}

func (p *Pg) GetReviewEvent(ctx context.Context, ownerID uuid.UUID, eventID uuid.UUID) (*models.ReviewEvent, error) {
	var reviewEvent models.ReviewEvent
	err := p.db.GetContext(ctx, &reviewEvent, "SELECT * FROM reviews_events WHERE owner_id = $1 AND event_id = $2 AND is_deleted = false", ownerID, eventID)

	return &reviewEvent, err
}

func (p *Pg) GetReviewsEvent(ctx context.Context, eventID uuid.UUID) ([]models.ReviewEvent, error) {
	var reviewEvent []models.ReviewEvent
	err := p.db.SelectContext(ctx, &reviewEvent, "SELECT * FROM reviews_events WHERE event_id = $1 AND is_deleted = false", eventID)

	return reviewEvent, err
}
