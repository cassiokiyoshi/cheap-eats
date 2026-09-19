package handlers

import (
	"fmt"
	"strings"
)

func normalizeOptionalName(value *string, field string) (*string, error) {
	if value == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, fmt.Errorf("%s must not be blank", field)
	}

	return &trimmed, nil
}
