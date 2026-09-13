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

INSERT INTO dishes (
    restaurant_id,
    name,
    price,
    currency
)
SELECT
    restaurants.id,
    'Shoyu Ramen',
    850,
    'JPY'
FROM restaurants
WHERE restaurants.name = 'Tokyo Ramen'
  AND restaurants.address = 'Marunouchi, Tokyo'
  AND NOT EXISTS (
      SELECT 1
      FROM dishes
      WHERE dishes.restaurant_id = restaurants.id
        AND dishes.name = 'Shoyu Ramen'
  );

INSERT INTO dishes (
    restaurant_id,
    name,
    price,
    currency
)
SELECT
    restaurants.id,
    'Gyudon',
    650,
    'JPY'
FROM restaurants
WHERE restaurants.name = 'Cheap Bowl'
  AND restaurants.address = 'Kanda, Tokyo'
  AND NOT EXISTS (
      SELECT 1
      FROM dishes
      WHERE dishes.restaurant_id = restaurants.id
        AND dishes.name = 'Gyudon'
  );

COMMIT;
