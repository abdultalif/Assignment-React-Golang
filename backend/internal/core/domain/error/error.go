package error

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrOutOfStock      = errors.New("out of stock")
)
