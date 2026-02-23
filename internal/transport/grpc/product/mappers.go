package product

import (
	"math/big"
	"strconv"
	"time"

	productv1 "product-catalog-service/proto/product/v1"
)

// RatToString formats *big.Rat as decimal string (e.g. "19.99").
func RatToString(r *big.Rat) string {
	if r == nil {
		return "0"
	}
	return r.FloatString(2)
}

func ProtoToCreateRequest(req *productv1.CreateProductRequest) CreateProductRequest {
	num, denom := req.BasePriceNumerator, req.BasePriceDenominator
	if denom == 0 {
		denom = 1
	}
	return CreateProductRequest{
		Name:             req.Name,
		Description:      req.Description,
		Category:         req.Category,
		BasePriceNum:    num,
		BasePriceDenom:  denom,
	}
}

func ProtoToUpdateRequest(req *productv1.UpdateProductRequest) UpdateProductRequest {
	return UpdateProductRequest{
		ProductID:   req.ProductId,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
	}
}

func ProtoToActivateRequest(req *productv1.ActivateProductRequest) ActivateProductRequest {
	return ActivateProductRequest{ProductID: req.ProductId}
}

func ProtoToDeactivateRequest(req *productv1.DeactivateProductRequest) DeactivateProductRequest {
	return DeactivateProductRequest{ProductID: req.ProductId}
}

func ProtoToApplyDiscountRequest(req *productv1.ApplyDiscountRequest) ApplyDiscountRequest {
	return ApplyDiscountRequest{
		ProductID:     req.ProductId,
		Percent:       req.Percent,
		StartDateUnix: req.StartDateUnix,
		EndDateUnix:   req.EndDateUnix,
	}
}

func ProtoToRemoveDiscountRequest(req *productv1.RemoveDiscountRequest) RemoveDiscountRequest {
	return RemoveDiscountRequest{ProductID: req.ProductId}
}

func ProtoToArchiveRequest(req *productv1.ArchiveProductRequest) ArchiveProductRequest {
	return ArchiveProductRequest{ProductID: req.ProductId}
}

func ProtoToListRequest(req *productv1.ListProductsRequest) ListProductsRequest {
	return ListProductsRequest{
		Category: req.Category,
		Status:   req.Status,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}
}

func DTOToGetReply(dto *GetProductDTO) *productv1.GetProductReply {
	if dto == nil {
		return nil
	}
	return &productv1.GetProductReply{
		ProductId:       dto.ProductID,
		Name:            dto.Name,
		Description:     dto.Description,
		Category:        dto.Category,
		BasePrice:       dto.BasePrice,
		EffectivePrice:  dto.EffectivePrice,
		DiscountPercent: dto.DiscountPercent,
		Status:          dto.Status,
	}
}

func DTOToListReply(result *ListProductsResultDTO) *productv1.ListProductsReply {
	if result == nil {
		return &productv1.ListProductsReply{}
	}
	out := make([]*productv1.ProductSummary, len(result.Products))
	for i, p := range result.Products {
		out[i] = &productv1.ProductSummary{
			ProductId:      p.ProductID,
			Name:           p.Name,
			Description:    p.Description,
			Category:       p.Category,
			BasePrice:      p.BasePrice,
			EffectivePrice: p.EffectivePrice,
			Status:         p.Status,
		}
	}
	return &productv1.ListProductsReply{Products: out, Total: result.Total}
}

// GetProductDTO is the query result for GetProduct reply mapping.
type GetProductDTO struct {
	ProductID       string
	Name            string
	Description     string
	Category        string
	BasePrice       string
	EffectivePrice  string
	DiscountPercent *int64
	Status          string
}

// ProductSummaryDTO is a list item for ListProducts reply.
type ProductSummaryDTO struct {
	ProductID       string
	Name            string
	Description     string
	Category        string
	BasePrice       string
	EffectivePrice  string
	Status          string
}

// CreateProductRequest for usecase (internal).
type CreateProductRequest struct {
	Name             string
	Description      string
	Category         string
	BasePriceNum     int64
	BasePriceDenom   int64
}

type UpdateProductRequest struct {
	ProductID   string
	Name        string
	Description string
	Category    string
}

type ActivateProductRequest struct {
	ProductID string
}

type DeactivateProductRequest struct {
	ProductID string
}

type ApplyDiscountRequest struct {
	ProductID     string
	Percent       int64
	StartDateUnix int64
	EndDateUnix   int64
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
	Limit    int32
	Offset   int32
}

// UnixToTime converts unix seconds to time.Time.
func UnixToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

// ParseRat parses a decimal string to *big.Rat (for tests).
func ParseRat(s string) *big.Rat {
	r, _ := new(big.Rat).SetString(s)
	return r
}

func _() { _ = strconv.Atoi } // use strconv if needed
