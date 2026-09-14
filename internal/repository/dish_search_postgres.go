package repository

import (
	"context"
	"fmt"
	"math"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
)

func (r *PostgresDishRepository) SearchNearby(
	ctx context.Context,
	latitude float64,
	longitude float64,
	radiusMeters float64,
	maxPrice int,
	sortBy string,
) ([]models.DishSearchResult, error) {
	var orderBy string

	switch sortBy {
	case "", "distance":
		orderBy = "distance_meters ASC, d.price ASC, d.id ASC"
	case "price":
		orderBy = "d.price ASC, distance_meters ASC, d.id ASC"
	default:
		return nil, fmt.Errorf("unsupported sort option: %q", sortBy)
	}

	query := `
		SELECT
			d.id,
			d.restaurant_id,
			d.name,
			d.price,
			d.currency,
			r.id,
			r.name,
			r.address,
			ST_Y(r.location::geometry),
			ST_X(r.location::geometry),
			ST_Distance(
				r.location,
				ST_SetSRID(
					ST_MakePoint($2, $1),
					4326
				)::geography
			) AS distance_meters
		FROM dishes d
		JOIN restaurants r ON r.id = d.restaurant_id
		WHERE d.price <= $4
		  AND ST_DWithin(
			  r.location,
			  ST_SetSRID(
				  ST_MakePoint($2, $1),
				  4326
			  )::geography,
			  $3
		  )
		ORDER BY ` + orderBy

	rows, err := r.pool.Query(
		ctx,
		query,
		latitude,
		longitude,
		radiusMeters,
		maxPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("search nearby dishes: %w", err)
	}
	defer rows.Close()

	results := make([]models.DishSearchResult, 0)

	for rows.Next() {
		var result models.DishSearchResult
		var distance float64

		if err := rows.Scan(
			&result.Dish.ID,
			&result.Dish.RestaurantID,
			&result.Dish.Name,
			&result.Dish.Price,
			&result.Dish.Currency,
			&result.Restaurant.ID,
			&result.Restaurant.Name,
			&result.Restaurant.Address,
			&result.Restaurant.Latitude,
			&result.Restaurant.Longitude,
			&distance,
		); err != nil {
			return nil, fmt.Errorf("scan nearby dish: %w", err)
		}

		result.DistanceMeters = int(math.Round(distance))
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nearby dishes: %w", err)
	}

	return results, nil
}
