BEGIN;

INSERT INTO restaurants (
    name,
    address,
    location
)
SELECT
    'Tokyo Ramen',
    'Marunouchi, Tokyo',
    ST_SetSRID(
        ST_MakePoint(139.7671, 35.6812),
        4326
    )::geography
WHERE NOT EXISTS (
    SELECT 1
    FROM restaurants
    WHERE name = 'Tokyo Ramen'
      AND address = 'Marunouchi, Tokyo'
);

INSERT INTO restaurants (
    name,
    address,
    location
)
SELECT
    'Cheap Bowl',
    'Kanda, Tokyo',
    ST_SetSRID(
        ST_MakePoint(139.7709, 35.6917),
        4326
    )::geography
WHERE NOT EXISTS (
    SELECT 1
    FROM restaurants
    WHERE name = 'Cheap Bowl'
      AND address = 'Kanda, Tokyo'
);

COMMIT;
