package product

import (
	"context"

	"product-catalog-service/internal/app/product/domain"
	productv1 "product-catalog-service/proto/product/v1"
)

// Handler implements product.v1.ProductService gRPC server.
type Handler struct {
	productv1.UnimplementedProductServiceServer

	CreateProduct   CreateProductRunner
	UpdateProduct   UpdateProductRunner
	ActivateProduct ActivateProductRunner
	DeactivateProduct DeactivateProductRunner
	ApplyDiscount   ApplyDiscountRunner
	RemoveDiscount  RemoveDiscountRunner
	ArchiveProduct  ArchiveProductRunner
	GetProduct      GetProductRunner
	ListProducts    ListProductsRunner
}

// CreateProduct creates a new product.
func (h *Handler) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.CreateProductReply, error) {
	if req == nil || req.Name == "" || req.Category == "" {
		return nil, MapDomainErrorToGRPC(domain.ErrInvalidInput)
	}
	appReq := ProtoToCreateRequest(req)
	productID, err := h.CreateProduct.Execute(ctx, appReq)
	if err != nil {
		return nil, MapDomainErrorToGRPC(err)
	}
	return &productv1.CreateProductReply{ProductId: productID}, nil
}

// UpdateProduct updates product details.
func (h *Handler) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.UpdateProductReply, error) {
	if req == nil || req.ProductId == "" {
		return nil, MapDomainErrorToGRPC(domain.ErrInvalidInput)
	}
	appReq := ProtoToUpdateRequest(req)
	if err := h.UpdateProduct.Execute(ctx, appReq); err != nil {
		return nil, MapDomainErrorToGRPC(err)
	}
	return &productv1.UpdateProductReply{}, nil
}

// ActivateProduct activates a product.
func (h *Handler) ActivateProduct(ctx context.Context, req *productv1.ActivateProductRequest) (*productv1.ActivateProductReply, error) {
	if req == nil || req.ProductId == "" {
		return nil, MapDomainErrorToGRPC(domain.ErrInvalidInput)
	}
	if err := h.ActivateProduct.Execute(ctx, ProtoToActivateRequest(req)); err != nil {
		return nil, MapDomainErrorToGRPC(err)
	}
	return &productv1.ActivateProductReply{}, nil
}

// DeactivateProduct deactivates a product.
func (h *Handler) DeactivateProduct(ctx context.Context, req *productv1.DeactivateProductRequest) (*productv1.DeactivateProductReply, error) {
	if req == nil || req.ProductId == "" {
		return nil, MapDomainErrorToGRPC(domain.ErrInvalidInput)
	}
	if err := h.DeactivateProduct.Execute(ctx, ProtoToDeactivateRequest(req)); err != nil {
		return nil, MapDomainErrorToGRPC(err)
	}
	return &productv1.DeactivateProductReply{}, nil
}

// ApplyDiscount applies a percentage discount.
func (h *Handler) ApplyDiscount(ctx context.Context, req *productv1.ApplyDiscountRequest) (*productv1.ApplyDiscountReply, error) {
	if req == nil || req.ProductId == "" {
		return nil, MapDomainErrorToGRPC(domain.ErrInvalidInput)
	}
	appReq := ProtoToApplyDiscountRequest(req)
	if err := h.ApplyDiscount.Execute(ctx, appReq); err != nil {
		return nil, MapDomainErrorToGRPC(err)
	}
	return &productv1.ApplyDiscountReply{}, nil
}

// RemoveDiscount removes the current discount.
func (h *Handler) RemoveDiscount(ctx context.Context, req *productv1.RemoveDiscountRequest) (*productv1.RemoveDiscountReply, error) {
	if req == nil || req.ProductId == "" {
		return nil, MapDomainErrorToGRPC(domain.ErrInvalidInput)
	}
	if err := h.RemoveDiscount.Execute(ctx, ProtoToRemoveDiscountRequest(req)); err != nil {
		return nil, MapDomainErrorToGRPC(err)
	}
	return &productv1.RemoveDiscountReply{}, nil
}

// ArchiveProduct soft-deletes a product.
func (h *Handler) ArchiveProduct(ctx context.Context, req *productv1.ArchiveProductRequest) (*productv1.ArchiveProductReply, error) {
	if req == nil || req.ProductId == "" {
		return nil, MapDomainErrorToGRPC(domain.ErrInvalidInput)
	}
	if err := h.ArchiveProduct.Execute(ctx, ProtoToArchiveRequest(req)); err != nil {
		return nil, MapDomainErrorToGRPC(err)
	}
	return &productv1.ArchiveProductReply{}, nil
}

// GetProduct returns a product by ID with effective price.
func (h *Handler) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductReply, error) {
	if req == nil || req.ProductId == "" {
		return nil, MapDomainErrorToGRPC(domain.ErrInvalidInput)
	}
	dto, err := h.GetProduct.Execute(ctx, req.ProductId)
	if err != nil {
		return nil, MapDomainErrorToGRPC(err)
	}
	return DTOToGetReply(dto), nil
}

// ListProducts returns paginated products.
func (h *Handler) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsReply, error) {
	appReq := ProtoToListRequest(req)
	result, err := h.ListProducts.Execute(ctx, appReq)
	if err != nil {
		return nil, MapDomainErrorToGRPC(err)
	}
	return DTOToListReply(result), nil
}
