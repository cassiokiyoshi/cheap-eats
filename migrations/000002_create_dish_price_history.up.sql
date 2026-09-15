CREATE TABLE dish_price_history (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dish_id BIGINT NOT NULL
        REFERENCES dishes(id)
        ON DELETE CASCADE,
    old_price INTEGER NOT NULL CHECK (old_price > 0),
    new_price INTEGER NOT NULL CHECK (new_price > 0),
    currency VARCHAR(3) NOT NULL
        CHECK (
            char_length(currency) = 3
            AND currency = upper(currency)
        ),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (old_price <> new_price)
);

CREATE INDEX dish_price_history_dish_id_id_idx
    ON dish_price_history (dish_id, id DESC);
