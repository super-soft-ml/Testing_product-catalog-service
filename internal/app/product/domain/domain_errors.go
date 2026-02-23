package domain

import "errors"

// Sentinel domain errors for product operations.
var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductNotActive      = errors.New("product is not active")
	ErrProductAlreadyActive  = errors.New("product is already active")
	ErrProductAlreadyArchived = errors.New("product is already archived")
	ErrInvalidDiscountPeriod = errors.New("invalid discount period")
	ErrDiscountAlreadyActive = errors.New("only one active discount per product allowed")
	ErrInvalidPrice          = errors.New("invalid price")
	ErrInvalidInput          = errors.New("invalid input")
)
