-- +goose Up

-- Основные таблицы
    CREATE TABLE IF NOT EXISTS users (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        vk_id BIGINT NOT NULL,
        passed_app_onboarding BOOL DEFAULT FALSE,
        passed_prisma_onboarding BOOL DEFAULT FALSE
    );
    CREATE INDEX idx_users_vkid ON users (vk_id);

    CREATE TABLE IF NOT EXISTS companies (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        user_id UUID NOT NULL,
        is_organisation BOOL NOT NULL,
        name VARCHAR(100) NOT NULL,
        description TEXT NOT NULL,
        average_rating DOUBLE PRECISION NOT NULL,
        photo_card TEXT NOT NULL
    );
    CREATE INDEX idx_companies_user_id ON companies (user_id);
    CREATE INDEX idx_companies_name ON companies (name);

    CREATE TABLE IF NOT EXISTS achievements (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        description TEXT NOT NULL,
        icon TEXT NOT NULL,
        coins INT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS routes (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        description TEXT NOT NULL,
        address_text TEXT NOT NULL,
        address_lng DOUBLE PRECISION NOT NULL,
        address_lat DOUBLE PRECISION NOT NULL,
        is_deleted BOOL NOT NULL
    );
    CREATE INDEX idx_routes_name ON routes (name);

    CREATE TABLE IF NOT EXISTS events (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        company_id UUID,
        name VARCHAR(100) NOT NULL,
        description TEXT NOT NULL,
        carousel TEXT[],
        icon TEXT,
        start_time TIMESTAMPTZ NOT NULL,
        address_text TEXT NOT NULL,
        address_lng DOUBLE PRECISION NOT NULL,
        address_lat DOUBLE PRECISION NOT NULL,
        is_deleted BOOL NOT NULL
    );
    CREATE INDEX idx_events_company_id ON events (company_id);
    CREATE INDEX idx_events_name ON events (name);

    CREATE TABLE IF NOT EXISTS places (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        description TEXT NOT NULL,
        carousel TEXT[],
        address_text TEXT NOT NULL,
        address_lng DOUBLE PRECISION NOT NULL,
        address_lat DOUBLE PRECISION NOT NULL,
        is_deleted BOOL NOT NULL
    );
    CREATE INDEX idx_places_name ON places (name);

    CREATE TABLE IF NOT EXISTS reviews_places (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        owner_id UUID NOT NULL,
        place_id UUID NOT NULL,
        review_text TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL,
        is_deleted BOOL NOT NULL
    );
    CREATE INDEX idx_reviews_places_owner ON reviews_places (owner_id);
    CREATE INDEX idx_reviews_places_route ON reviews_places (place_id);

    CREATE TABLE IF NOT EXISTS reviews_routes (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        owner_id UUID NOT NULL,
        route_id UUID NOT NULL,
        review_text TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL,
        is_deleted BOOL NOT NULL
    );
    CREATE INDEX idx_reviews_routes_owner ON reviews_routes (owner_id);
    CREATE INDEX idx_reviews_routes_route ON reviews_routes (route_id);

    CREATE TABLE IF NOT EXISTS reviews_events (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        owner_id UUID NOT NULL,
        event_id UUID NOT NULL,
        review_text TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL,
        is_deleted BOOL NOT NULL
    );
    CREATE INDEX idx_reviews_events_owner ON reviews_events (owner_id);
    CREATE INDEX idx_reviews_events_route ON reviews_events (event_id);


-- Фильтры
    CREATE TABLE IF NOT EXISTS events_filters (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        name VARCHAR(100)
    );

    CREATE TABLE IF NOT EXISTS map_filters (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        name VARCHAR(100)
    );


-- Связи между таблицами
    CREATE TABLE IF NOT EXISTS users_map_filters_rel (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        user_id UUID,
        filter_id UUID
    );
    CREATE INDEX idx_users_map_filters_rel_user_filter ON users_map_filters_rel (user_id, filter_id);

    CREATE TABLE IF NOT EXISTS users_privacy_rel (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        user_id UUID,
        can_view_achievements BOOL,
        can_view_progress_on_map BOOL
    );
    CREATE INDEX idx_users_privacy_rel_user ON users_privacy_rel (user_id);

    CREATE TABLE IF NOT EXISTS users_progress_on_map_rel (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        user_id UUID,
        route_id UUID,
        event_id UUID,
        place_id UUID,
        created_at TIMESTAMPTZ
    );
    CREATE INDEX idx_users_progress_on_map_rel_user ON users_progress_on_map_rel (user_id);

    CREATE TABLE IF NOT EXISTS users_achievements_rel (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        user_id UUID NOT NULL,
        achievement_id UUID NOT NULL
    );
    CREATE INDEX idx_users_achievements_rel_user_rel ON users_achievements_rel (user_id);

    CREATE TABLE IF NOT EXISTS users_coins_rel (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        user_id UUID NOT NULL,
        coins INT NOT NULL DEFAULT 0,
        operation BOOL NOT NULL
    );
    CREATE INDEX idx_users_coins_rel_user ON users_coins_rel (user_id);

    CREATE TABLE IF NOT EXISTS events_filters_rel (
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        event_id UUID NOT NULL,
        filter_id UUID NOT NULL
    );
    CREATE INDEX idx_events_filters_rel_event_filter ON events_filters_rel (event_id, filter_id);

-- +goose Down
