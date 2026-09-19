package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDishRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresDishRepository(
	pool *pgxpool.Pool,
) *PostgresDishRepository {
	return &PostgresDishRepository{
		pool: pool,
	}
}

func (r *PostgresDishRepository) List(
	ctx context.Context,
) ([]models.Dish, error) {
	const query = `
		SELECT id, restaurant_id, name, name_ja, name_en, price, currency
		FROM dishes
		ORDER BY id
	`

	return r.queryDishes(ctx, query)
}

func (r *PostgresDishRepository) ListByMaxPrice(
	ctx context.Context,
	maxPrice int,
) ([]models.Dish, error) {
	const query = `
		SELECT id, restaurant_id, name, name_ja, name_en, price, currency
		FROM dishes
		WHERE price <= $1
		ORDER BY price, id
	`

	return r.queryDishes(ctx, query, maxPrice)
}

func (r *PostgresDishRepository) FindByID(
	ctx context.Context,
	id int64,
) (models.Dish, bool, error) {
	const query = `
		SELECT id, restaurant_id, name, name_ja, name_en, price, currency
		FROM dishes
		WHERE id = $1
	`

	var dish models.Dish

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&dish.ID,
		&dish.RestaurantID,
		&dish.Name,
		&dish.NameJA,
		&dish.NameEN,
		&dish.Price,
		&dish.Currency,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Dish{}, false, nil
	}
	if err != nil {
		return models.Dish{}, false, fmt.Errorf(
			"find dish by ID: %w",
			err,
		)
	}

	return dish, true, nil
}

func (r *PostgresDishRepository) Create(
	ctx context.Context,
	dish models.Dish,
) (models.Dish, error) {
	const query = `
		INSERT INTO dishes (
			restaurant_id,
			name,
			name_ja,
			name_en,
			price,
			currency
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		dish.RestaurantID,
		dish.Name,
		dish.NameJA,
		dish.NameEN,
		dish.Price,
		dish.Currency,
	).Scan(&dish.ID)
	if err != nil {
		return models.Dish{}, fmt.Errorf(
			"create dish: %w",
			err,
		)
	}

	return dish, nil
}

func (r *PostgresDishRepository) queryDishes(
	ctx context.Context,
	query string,
	arguments ...any,
) ([]models.Dish, error) {
	rows, err := r.pool.Query(ctx, query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("query dishes: %w", err)
	}
	defer rows.Close()

	dishes := make([]models.Dish, 0)

	for rows.Next() {
		var dish models.Dish

		if err := rows.Scan(
			&dish.ID,
			&dish.RestaurantID,
			&dish.Name,
			&dish.NameJA,
			&dish.NameEN,
			&dish.Price,
			&dish.Currency,
		); err != nil {
			return nil, fmt.Errorf("scan dish: %w", err)
		}

		dishes = append(dishes, dish)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dishes: %w", err)
	}

	return dishes, nil
}

func (r *PostgresDishRepository) ListByRestaurantID(
	ctx context.Context,
	restaurantID int64,
) ([]models.Dish, error) {
	const query = `
		SELECT id, restaurant_id, name, name_ja, name_en, price, currency
		FROM dishes
		WHERE restaurant_id = $1
		ORDER BY id
	`

	return r.queryDishes(ctx, query, restaurantID)
}

func (r *PostgresDishRepository) UpdatePrice(
	ctx context.Context,
	id int64,
	price int,
) (models.Dish, bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.Dish{}, false, fmt.Errorf(
			"begin price update: %w", err,
		)
	}
	defer tx.Rollback(ctx)

	const findQuery = `
		SELECT id, restaurant_id, name, name_ja, name_en, price, currency
		FROM dishes
		WHERE id = $1
		FOR UPDATE
	`

	var dish models.Dish

	err = tx.QueryRow(ctx, findQuery, id).Scan(
		&dish.ID,
		&dish.RestaurantID,
		&dish.Name,
		&dish.NameJA,
		&dish.NameEN,
		&dish.Price,
		&dish.Currency,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Dish{}, false, nil
	}
	if err != nil {
		return models.Dish{}, false, fmt.Errorf(
			"lock dish for price update: %w", err,
		)
	}

	// An unchanged price is successful but creates no history.
	if dish.Price == price {
		return dish, true, nil
	}

	const updateQuery = `
		UPDATE dishes
		SET price = $2, updated_at = NOW()
		WHERE id = $1
	`

	_, err = tx.Exec(ctx, updateQuery, id, price)
	if err != nil {
		return models.Dish{}, false, fmt.Errorf(
			"update dish price: %w", err,
		)
	}

	const historyQuery = `
		INSERT INTO dish_price_history (
			dish_id, old_price, new_price, currency
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(
		ctx,
		historyQuery,
		id,
		dish.Price,
		price,
		dish.Currency,
	)
	if err != nil {
		return models.Dish{}, false, fmt.Errorf(
			"record dish price history: %w", err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Dish{}, false, fmt.Errorf(
			"commit price update: %w", err,
		)
	}

	dish.Price = price
	return dish, true, nil
}

var _ DishStore = (*PostgresDishRepository)(nil)
