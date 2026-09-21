BEGIN;

-- Protect against accidentally seeding the integration-test database.
DO $$
BEGIN
    IF current_database() <> 'cheap_eats' THEN
        RAISE EXCEPTION
            'Demo seed must run against cheap_eats, not %',
            current_database();
    END IF;
END
$$;

CREATE TEMP TABLE demo_restaurants (
    code TEXT PRIMARY KEY,
    name_ja TEXT NOT NULL,
    name_en TEXT NOT NULL,
    area TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL
) ON COMMIT DROP;

INSERT INTO demo_restaurants VALUES
    ('komorebi', 'こもれび食堂', 'Komorebi Kitchen',
     'Marunouchi', 35.6815, 139.7655),
    ('tsubame', 'つばめ麺処', 'Tsubame Noodles',
     'Marunouchi', 35.6803, 139.7660),
    ('hinata', 'ひなたカレー', 'Hinata Curry',
     'Yaesu', 35.6814, 139.7690),
    ('marufuku', 'まる福ごはん', 'Marufuku Rice Kitchen',
     'Yaesu', 35.6800, 139.7687),
    ('nagi', 'なぎそば', 'Nagi Soba',
     'Marunouchi', 35.6830, 139.7660),
    ('koharu', 'こはる屋', 'Koharu Kitchen',
     'Yaesu', 35.6825, 139.7692),
    ('akari', 'あかり食堂', 'Akari Diner',
     'Kyobashi', 35.6768, 139.7700),
    ('suzume', 'すずめ亭', 'Suzume Kitchen',
     'Nihonbashi', 35.6840, 139.7740),
    ('aoba', '青葉めし', 'Aoba Rice Kitchen',
     'Kanda', 35.6900, 139.7700),
    ('yuzuki', 'ゆづき食堂', 'Yuzuki Diner',
     'Kanda', 35.6920, 139.7680);

INSERT INTO restaurants (
    name, name_ja, name_en, address, location
)
SELECT
    '[DEMO] ' || demo.code,
    demo.name_ja,
    demo.name_en,
    demo.area || ', Tokyo — fictional demo location',
    ST_SetSRID(
        ST_MakePoint(demo.longitude, demo.latitude),
        4326
    )::geography
FROM demo_restaurants demo
WHERE NOT EXISTS (
    SELECT 1
    FROM restaurants existing
    WHERE existing.name = '[DEMO] ' || demo.code
);

CREATE TEMP TABLE demo_dishes (
    restaurant_code TEXT NOT NULL,
    name_ja TEXT NOT NULL,
    name_en TEXT NOT NULL,
    price INTEGER NOT NULL
) ON COMMIT DROP;

INSERT INTO demo_dishes VALUES
    ('komorebi', '唐揚げ定食', 'Karaage Set', 850),
    ('komorebi', '焼き鯖定食', 'Grilled Mackerel Set', 900),
    ('tsubame', '醤油ラーメン', 'Shoyu Ramen', 780),
    ('tsubame', '餃子', 'Gyoza', 380),
    ('hinata', 'カレーライス', 'Japanese Curry', 650),
    ('hinata', '牛丼', 'Gyudon', 600),
    ('marufuku', '牛丼', 'Gyudon', 580),
    ('marufuku', '親子丼定食', 'Oyakodon Set', 800),
    ('nagi', 'ざるそば', 'Zaru Soba', 600),
    ('nagi', 'かけうどん', 'Udon', 500),
    ('koharu', 'たこ焼き', 'Takoyaki', 500),
    ('koharu', '餃子', 'Gyoza', 400),
    ('akari', '唐揚げ定食', 'Karaage Set', 880),
    ('akari', 'カレーライス', 'Japanese Curry', 700),
    ('suzume', '醤油ラーメン', 'Shoyu Ramen', 820),
    ('suzume', '餃子', 'Gyoza', 420),
    ('aoba', '親子丼定食', 'Oyakodon Set', 750),
    ('aoba', '牛丼', 'Gyudon', 550),
    ('yuzuki', '焼き鯖定食', 'Grilled Mackerel Set', 950),
    ('yuzuki', 'ざるそば', 'Zaru Soba', 650);

INSERT INTO dishes (
    restaurant_id, name, name_ja, name_en, price, currency
)
SELECT
    restaurant.id,
    demo.name_en,
    demo.name_ja,
    demo.name_en,
    demo.price,
    'JPY'
FROM demo_dishes demo
JOIN restaurants restaurant
    ON restaurant.name = '[DEMO] ' || demo.restaurant_code
WHERE NOT EXISTS (
    SELECT 1
    FROM dishes existing
    WHERE existing.restaurant_id = restaurant.id
      AND existing.name = demo.name_en
);

COMMIT;
