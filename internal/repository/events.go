package repository

import (
	"context"

	"github.com/ShpullRequest/backend/internal/models"
	"github.com/google/uuid"
)

func (p *Pg) NewEvent(ctx context.Context, event models.Event) (*models.Event, error) {
	_, err := p.db.ExecContext(
		ctx,
		"INSERT INTO events (company_id, name, description, carousel, icon, start_time, address_text, address_lng, address_lat, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		event.CompanyID,
		event.Name,
		event.Description,
		event.Carousel,
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

	return &event, nil
}

func (p *Pg) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	var event models.Event
	err := p.db.GetContext(ctx, &event, "SELECT * FROM events WHERE id = $1", id)

	return &event, err
}

func (p *Pg) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	err := p.db.SelectContext(ctx, &events, "SELECT * FROM events")

	return events, err
}

func (p *Pg) GetEventByCompanyID(ctx context.Context, eventID uuid.UUID, companyID uuid.UUID) (*models.Event, error) {
	var event models.Event
	err := p.db.GetContext(ctx, &event, "SELECT * FROM events WHERE eventID = $1 AND company_id = $2", eventID, companyID)

	return &event, err
}

func (p *Pg) GetAllEventsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]models.Event, error) {
	var events []models.Event
	err := p.db.SelectContext(ctx, &events, "SELECT * FROM events WHERE company_id = $1", companyID)

	return events, err
}

func (p *Pg) NewReviewEvent(ctx context.Context, reviewEvent models.ReviewEvent) (*models.ReviewEvent, error) {
	_, err := p.db.ExecContext(
		ctx,
		"INSERT INTO reviews_events (owner_id, event_id, review_text, created_at, is_deleted) VALUES ($1, $2, $3, $4, $5)",
		reviewEvent.OwnerID,
		reviewEvent.EventID,
		reviewEvent.ReviewText,
		reviewEvent.CreatedAt,
		reviewEvent.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &reviewEvent, nil
}

func (p *Pg) GetReviewsEvent(ctx context.Context, eventID uuid.UUID) ([]models.ReviewEvent, error) {
	var reviewEvent []models.ReviewEvent
	err := p.db.SelectContext(ctx, &reviewEvent, "SELECT * FROM reviews_events WHERE event_id = $1", eventID)

	return reviewEvent, err
}
