package handlers

import "errors"

func validatePrice(price int) error {
	// PostgreSQL INTEGER is a signed 32-bit integer.
	const maxPrice = 2_147_483_647

	if price < 1 || price > maxPrice {
		return errors.New("price must be between 1 and 2147483647")
	}

	return nil
}
