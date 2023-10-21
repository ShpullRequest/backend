-- +goose Up

-- +goose StatementBegin
    CREATE OR REPLACE FUNCTION calculate_company_rating(c_id UUID)
        RETURNS DOUBLE PRECISION AS $$
    BEGIN
        RETURN COALESCE(AVG(stars), 0) FROM (
            SELECT stars FROM reviews_routes WHERE route_id IN (SELECT id FROM routes WHERE company_id = c_id) AND is_deleted = false
            UNION ALL
            SELECT stars FROM reviews_events WHERE event_id IN (SELECT id FROM events WHERE company_id = c_id) AND is_deleted = false
        ) AS rating;
    END;
    $$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down