package error

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrOutOfStock      = errors.New("out of stock")

	ErrOrderNotFound = errors.New("order not found")

	ErrInvalidDateRange = errors.New("invalid date range")

	ErrJobNotFound = errors.New("job not found")
)
