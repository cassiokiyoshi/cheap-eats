BEGIN;

ALTER TABLE restaurants
    DROP COLUMN name_en,
    DROP COLUMN name_ja;

ALTER TABLE dishes
    DROP COLUMN name_en,
    DROP COLUMN name_ja;

COMMIT;
