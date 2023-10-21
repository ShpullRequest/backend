-- +goose Up

-- Main tables
    CREATE TABLE IF NOT EXISTS users (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        VkId BIGINT NOT NULL,
        PassedAppOnboarding BOOL DEFAULT FALSE,
        PassedPrismaOnboarding BOOL DEFAULT FALSE
    );
    CREATE INDEX idx_users_vkid ON users (VkId);

    CREATE TABLE IF NOT EXISTS companies (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        UserId UUID NOT NULL,
        IsOrganisation BOOL NOT NULL,
        Name VARCHAR(100) NOT NULL,
        Description TEXT NOT NULL,
        AverageRating DOUBLE PRECISION NOT NULL,
        PhotoCard TEXT NOT NULL
    );
    CREATE INDEX idx_companies_user_id ON companies (UserId);
    CREATE INDEX idx_companies_name ON companies (Name);

    CREATE TABLE IF NOT EXISTS achievements (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        Name VARCHAR(100) NOT NULL,
        Description TEXT NOT NULL,
        Icon TEXT NOT NULL,
        Coins INT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS routes (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        Name VARCHAR(100) NOT NULL,
        Description TEXT NOT NULL,
        AddressText TEXT NOT NULL,
        AddressLng DOUBLE PRECISION NOT NULL,
        AddressLat DOUBLE PRECISION NOT NULL,
        IsDeleted BOOL NOT NULL
    );
    CREATE INDEX idx_routes_name ON routes (Name);

    CREATE TABLE IF NOT EXISTS events (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        CompanyId UUID,
        Name VARCHAR(100) NOT NULL,
        Description TEXT NOT NULL,
        Carousel TEXT[],
        Icon TEXT,
        StartTime TIMESTAMPTZ NOT NULL,
        AddressText TEXT NOT NULL,
        AddressLng DOUBLE PRECISION NOT NULL,
        AddressLat DOUBLE PRECISION NOT NULL,
        IsDeleted BOOL NOT NULL
    );
    CREATE INDEX idx_events_company_id ON events (CompanyId);
    CREATE INDEX idx_events_name ON events (Name);

    CREATE TABLE IF NOT EXISTS places (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        Name VARCHAR(100) NOT NULL,
        Description TEXT NOT NULL,
        Carousel TEXT[],
        AddressText TEXT NOT NULL,
        AddressLng DOUBLE PRECISION NOT NULL,
        AddressLat DOUBLE PRECISION NOT NULL,
        IsDeleted BOOL NOT NULL
    );
    CREATE INDEX idx_places_name ON places (Name);

    CREATE TABLE IF NOT EXISTS reviews_places (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        OwnerId UUID NOT NULL,
        PlaceId UUID NOT NULL,
        ReviewText TEXT NOT NULL,
        CreatedAt TIMESTAMPTZ NOT NULL,
        IsDeleted BOOL NOT NULL
    );
    CREATE INDEX idx_reviews_places_owner ON reviews_places (OwnerId);
    CREATE INDEX idx_reviews_places_route ON reviews_places (PlaceId);

    CREATE TABLE IF NOT EXISTS reviews_routes (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        OwnerId UUID NOT NULL,
        RouteId UUID NOT NULL,
        ReviewText TEXT NOT NULL,
        CreatedAt TIMESTAMPTZ NOT NULL,
        IsDeleted BOOL NOT NULL
    );
    CREATE INDEX idx_reviews_routes_owner ON reviews_routes (OwnerId);
    CREATE INDEX idx_reviews_routes_route ON reviews_routes (RouteId);

    CREATE TABLE IF NOT EXISTS reviews_events (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        OwnerId UUID NOT NULL,
        EventId UUID NOT NULL,
        ReviewText TEXT NOT NULL,
        CreatedAt TIMESTAMPTZ NOT NULL,
        IsDeleted BOOL NOT NULL
    );
    CREATE INDEX idx_reviews_events_owner ON reviews_events (OwnerId);
    CREATE INDEX idx_reviews_events_route ON reviews_events (EventId);


-- Filters
    CREATE TABLE IF NOT EXISTS events_filters (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        Name VARCHAR(100)
    );

    CREATE TABLE IF NOT EXISTS map_filters (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        Name VARCHAR(100)
    );


-- Tables relationship
    CREATE TABLE IF NOT EXISTS users_map_filters_rel (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        UserId UUID,
        FilterId UUID
    );
    CREATE INDEX idx_users_map_filters_rel_user_filter ON users_map_filters_rel (UserId, FilterId);

    CREATE TABLE IF NOT EXISTS users_privacy_rel (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        UserId UUID,
        CanViewAchievements BOOL,
        CanViewProgressOnMap BOOL
    );
    CREATE INDEX idx_users_privacy_rel_user ON users_privacy_rel (UserId);

    CREATE TABLE IF NOT EXISTS users_progress_on_map_rel (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        UserId UUID,
        RouteId UUID,
        EventId UUID,
        PlaceId UUID,
        CreatedAt TIMESTAMPTZ
    );
    CREATE INDEX idx_users_progress_on_map_rel_user ON users_progress_on_map_rel (UserId);

    CREATE TABLE IF NOT EXISTS users_achievements_rel (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        UserId UUID NOT NULL,
        AchievementId UUID NOT NULL
    );
    CREATE INDEX idx_users_achievements_rel_user_rel ON users_achievements_rel (UserId);

    CREATE TABLE IF NOT EXISTS users_coins_rel (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        UserId UUID NOT NULL,
        Coins INT NOT NULL DEFAULT 0,
        Operation BOOL NOT NULL
    );
    CREATE INDEX idx_users_coins_rel_user ON users_coins_rel (UserId);

    CREATE TABLE IF NOT EXISTS events_filters_rel (
        Id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        EventId UUID NOT NULL,
        FilterId UUID NOT NULL
    );
    CREATE INDEX idx_events_filters_rel_event_filter ON events_filters_rel (EventId, FilterId);

-- +goose Down