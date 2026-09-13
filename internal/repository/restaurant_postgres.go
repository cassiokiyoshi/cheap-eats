package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRestaurantRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRestaurantRepository(
	pool *pgxpool.Pool,
) *PostgresRestaurantRepository {
	return &PostgresRestaurantRepository{
		pool: pool,
	}
}

func (r *PostgresRestaurantRepository) List(
	ctx context.Context,
) ([]models.Restaurant, error) {
	const query = `
		SELECT
			id,
			name,
			address,
			ST_Y(location::geometry),
			ST_X(location::geometry)
		FROM restaurants
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list restaurants: %w", err)
	}
	defer rows.Close()

	restaurants := make([]models.Restaurant, 0)

	for rows.Next() {
		var restaurant models.Restaurant

		if err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.Address,
			&restaurant.Latitude,
			&restaurant.Longitude,
		); err != nil {
			return nil, fmt.Errorf(
				"scan restaurant: %w",
				err,
			)
		}

		restaurants = append(restaurants, restaurant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate restaurants: %w",
			err,
		)
	}

	return restaurants, nil
}

func (r *PostgresRestaurantRepository) FindByID(
	ctx context.Context,
	id int64,
) (models.Restaurant, bool, error) {
	const query = `
		SELECT
			id,
			name,
			address,
			ST_Y(location::geometry),
			ST_X(location::geometry)
		FROM restaurants
		WHERE id = $1
	`

	var restaurant models.Restaurant

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Address,
		&restaurant.Latitude,
		&restaurant.Longitude,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Restaurant{}, false, nil
	}
	if err != nil {
		return models.Restaurant{}, false, fmt.Errorf(
			"find restaurant by ID: %w",
			err,
		)
	}

	return restaurant, true, nil
}

func (r *PostgresRestaurantRepository) Nearby(
	ctx context.Context,
	latitude float64,
	longitude float64,
	radiusMeters float64,
) ([]models.RestaurantSuggestion, error) {
	const query = `
		SELECT
			id,
			name,
			address,
			ST_Y(location::geometry),
			ST_X(location::geometry),
			ROUND(
				ST_Distance(
					location,
					ST_SetSRID(
						ST_MakePoint($2, $1),
						4326
					)::geography
				)
			)::INTEGER AS distance_meters
		FROM restaurants
		WHERE ST_DWithin(
			location,
			ST_SetSRID(
				ST_MakePoint($2, $1),
				4326
			)::geography,
			$3
		)
		ORDER BY distance_meters
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		latitude,
		longitude,
		radiusMeters,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"search nearby restaurants: %w",
			err,
		)
	}
	defer rows.Close()

	suggestions := make([]models.RestaurantSuggestion, 0)

	for rows.Next() {
		var suggestion models.RestaurantSuggestion

		if err := rows.Scan(
			&suggestion.ID,
			&suggestion.Name,
			&suggestion.Address,
			&suggestion.Latitude,
			&suggestion.Longitude,
			&suggestion.DistanceMeters,
		); err != nil {
			return nil, fmt.Errorf(
				"scan nearby restaurant: %w",
				err,
			)
		}

		suggestions = append(suggestions, suggestion)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate nearby restaurants: %w",
			err,
		)
	}

	return suggestions, nil
}

var _ RestaurantStore = (*PostgresRestaurantRepository)(nil)
