package repository

import (
	"context"
	"fmt"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
)

func (r *PostgresDishRepository) ListPriceHistory(
	ctx context.Context,
	dishID int64,
	limit int,
	offset int,
) ([]models.DishPriceHistory, error) {
	if limit < 1 || limit > 100 || offset < 0 {
		return nil, fmt.Errorf("invalid price history pagination")
	}
	const query = `
		SELECT
			id,
			dish_id,
			old_price,
			new_price,
			currency,
			changed_at
		FROM dish_price_history
		WHERE dish_id = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, dishID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list dish price history: %w", err)
	}
	defer rows.Close()

	history := make([]models.DishPriceHistory, 0)

	for rows.Next() {
		var entry models.DishPriceHistory

		if err := rows.Scan(
			&entry.ID,
			&entry.DishID,
			&entry.OldPrice,
			&entry.NewPrice,
			&entry.Currency,
			&entry.ChangedAt,
		); err != nil {
			return nil, fmt.Errorf("scan dish price history: %w", err)
		}

		history = append(history, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dish price history: %w", err)
	}

	return history, nil
}
