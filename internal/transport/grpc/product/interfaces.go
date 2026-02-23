package product

import (
	"context"

	productv1 "product-catalog-service/proto/product/v1"
)

// CreateProductRunner runs the CreateProduct usecase.
type CreateProductRunner interface {
	Execute(ctx context.Context, req CreateProductRequest) (productID string, err error)
}

// UpdateProductRunner runs the UpdateProduct usecase.
type UpdateProductRunner interface {
	Execute(ctx context.Context, req UpdateProductRequest) error
}

// ActivateProductRunner runs the ActivateProduct usecase.
type ActivateProductRunner interface {
	Execute(ctx context.Context, req ActivateProductRequest) error
}

// DeactivateProductRunner runs the DeactivateProduct usecase.
type DeactivateProductRunner interface {
	Execute(ctx context.Context, req DeactivateProductRequest) error
}

// ApplyDiscountRunner runs the ApplyDiscount usecase.
type ApplyDiscountRunner interface {
	Execute(ctx context.Context, req ApplyDiscountRequest) error
}

// RemoveDiscountRunner runs the RemoveDiscount usecase.
type RemoveDiscountRunner interface {
	Execute(ctx context.Context, req RemoveDiscountRequest) error
}

// ArchiveProductRunner runs the ArchiveProduct usecase.
type ArchiveProductRunner interface {
	Execute(ctx context.Context, req ArchiveProductRequest) error
}

// GetProductRunner runs the GetProduct query.
type GetProductRunner interface {
	Execute(ctx context.Context, productID string) (*GetProductDTO, error)
}

// ListProductsRunner runs the ListProducts query.
type ListProductsRunner interface {
	Execute(ctx context.Context, req ListProductsRequest) (*ListProductsResultDTO, error)
}

// ListProductsResultDTO for reply mapping.
type ListProductsResultDTO struct {
	Products []*ProductSummaryDTO
	Total    int32
}
