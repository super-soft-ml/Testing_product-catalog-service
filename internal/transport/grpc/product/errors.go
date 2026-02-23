package product

import (
	"math/big"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"product-catalog-service/internal/app/product/domain"
)

// MapDomainErrorToGRPC maps domain errors to gRPC status codes.
func MapDomainErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case err == domain.ErrProductNotFound:
		return status.Error(codes.NotFound, err.Error())
	case err == domain.ErrProductNotActive, err == domain.ErrInvalidDiscountPeriod,
		err == domain.ErrProductAlreadyActive, err == domain.ErrProductAlreadyArchived,
		err == domain.ErrDiscountAlreadyActive:
		return status.Error(codes.FailedPrecondition, err.Error())
	case err == domain.ErrInvalidInput, err == domain.ErrInvalidPrice:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// RatToDecimalString formats *big.Rat as decimal string.
func RatToDecimalString(r *big.Rat) string {
	if r == nil {
		return "0"
	}
	f, _ := r.Float64()
	return strconv.FormatFloat(f, 'f', 2, 64)
}
