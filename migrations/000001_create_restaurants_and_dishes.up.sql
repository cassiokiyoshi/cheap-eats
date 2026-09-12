CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE restaurants (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    address TEXT NOT NULL CHECK (btrim(address) <> ''),
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX restaurants_location_idx
    ON restaurants
    USING GIST (location);

CREATE TABLE dishes (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_id BIGINT NOT NULL
        REFERENCES restaurants(id)
        ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    price INTEGER NOT NULL CHECK (price > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'JPY'
        CHECK (
            char_length(currency) = 3
            AND currency = upper(currency)
        ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX dishes_restaurant_id_idx
    ON dishes (restaurant_id);

CREATE INDEX dishes_price_idx
    ON dishes (price);
