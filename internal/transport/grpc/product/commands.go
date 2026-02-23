package product

import (
	"context"

	productv1 "product-catalog-service/proto/product/v1"
)

// Commands holds all command usecases for the handler.
type Commands struct {
	CreateProduct   CreateProductRunner
	UpdateProduct   UpdateProductRunner
	ActivateProduct ActivateProductRunner
	DeactivateProduct DeactivateProductRunner
	ApplyDiscount   ApplyDiscountRunner
	RemoveDiscount  RemoveDiscountRunner
	ArchiveProduct  ArchiveProductRunner
}

// Queries holds query usecases.
type Queries struct {
	GetProduct  GetProductRunner
	ListProducts ListProductsRunner
}

// CreateProductRunner runs CreateProduct usecase.
type CreateProductRunner interface {
	Execute(ctx context.Context, req CreateProductRequest) (string, error)
}

// UpdateProductRunner runs UpdateProduct usecase.
type UpdateProductRunner interface {
	Execute(ctx context.Context, req UpdateProductRequest) error
}

// ActivateProductRunner runs ActivateProduct usecase.
type ActivateProductRunner interface {
	Execute(ctx context.Context, req ActivateProductRequest) error
}

// DeactivateProductRunner runs DeactivateProduct usecase.
type DeactivateProductRunner interface {
	Execute(ctx context.Context, req ActivateProductRequest) error
}

// ApplyDiscountRunner runs ApplyDiscount usecase.
type ApplyDiscountRunner interface {
	Execute(ctx context.Context, req ApplyDiscountRequest) error
}

// RemoveDiscountRunner runs RemoveDiscount usecase.
type RemoveDiscountRunner interface {
	Execute(ctx context.Context, req RemoveDiscountRequest) error
}

// ArchiveProductRunner runs ArchiveProduct usecase.
type ArchiveProductRunner interface {
	Execute(ctx context.Context, req ArchiveProductRequest) error
}

// GetProductRunner runs GetProduct query.
type GetProductRunner interface {
	Execute(ctx context.Context, productID string) (*GetProductResult, error)
}

// ListProductsRunner runs ListProducts query.
type ListProductsRunner interface {
	Execute(ctx context.Context, req ListProductsRequest) (*ListProductsResult, error)
}

// App request/result types (mirror usecases).
type CreateProductRequest struct {
	Name, Description, Category string
	BasePriceNum, BasePriceDenom int64
}

type UpdateProductRequest struct {
	ProductID, Name, Description, Category string
}

type ActivateProductRequest struct {
	ProductID string
}

type ApplyDiscountRequest struct {
	ProductID string
	Percent   int64
	StartDateUnix, EndDateUnix int64
}

type RemoveDiscountRequest struct {
	ProductID string
}

type ArchiveProductRequest struct {
	ProductID string
}

type ListProductsRequest struct {
	Category string
	Status   string
	Limit, Offset int32
}

type GetProductResult struct {
	ProductID, Name, Description, Category string
	BasePrice, EffectivePrice              string
	DiscountPercent                        *int64
	Status                                 string
}

type ListProductsResult struct {
	Products []*ProductSummary
	Total    int32
}

type ProductSummary struct {
	ProductID, Name, Description, Category string
	BasePrice, EffectivePrice, Status       string
}
