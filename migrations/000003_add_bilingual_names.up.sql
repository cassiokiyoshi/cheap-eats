BEGIN;

ALTER TABLE dishes
    ADD COLUMN name_ja TEXT
        CHECK (name_ja IS NULL OR btrim(name_ja) <> ''),
    ADD COLUMN name_en TEXT
        CHECK (name_en IS NULL OR btrim(name_en) <> '');

ALTER TABLE restaurants
    ADD COLUMN name_ja TEXT
        CHECK (name_ja IS NULL OR btrim(name_ja) <> ''),
    ADD COLUMN name_en TEXT
        CHECK (name_en IS NULL OR btrim(name_en) <> '');

COMMIT;
