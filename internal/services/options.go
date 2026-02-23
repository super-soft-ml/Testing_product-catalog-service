package services

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"cloud.google.com/go/spanner"

	"product-catalog-service/internal/app/product/repo"
	"product-catalog-service/internal/app/product/usecases/activate_product"
	"product-catalog-service/internal/app/product/usecases/apply_discount"
	"product-catalog-service/internal/app/product/usecases/archive_product"
	"product-catalog-service/internal/app/product/usecases/create_product"
	"product-catalog-service/internal/app/product/usecases/deactivate_product"
	"product-catalog-service/internal/app/product/usecases/remove_discount"
	"product-catalog-service/internal/app/product/usecases/update_product"
	"product-catalog-service/internal/app/product/queries/get_product"
	"product-catalog-service/internal/app/product/queries/list_products"
	"product-catalog-service/internal/pkg/clock"
	"product-catalog-service/internal/pkg/committer"
	"product-catalog-service/internal/transport/grpc/product"
)

// Options holds all dependencies for the application.
type Options struct {
	SpannerClient *spanner.Client

	// Commands
	CreateProduct   *create_product.Interactor
	UpdateProduct   *update_product.Interactor
	ActivateProduct *activate_product.Interactor
	DeactivateProduct *deactivate_product.Interactor
	ApplyDiscount   *apply_discount.Interactor
	RemoveDiscount  *remove_discount.Interactor
	ArchiveProduct  *archive_product.Interactor

	// Queries
	GetProductQuery  get_product.Query
	ListProductsQuery list_products.Query
}

// NewOptions builds the DI container.
func NewOptions(client *spanner.Client) *Options {
	clk := clock.RealClock{}
	productRepo := repo.NewProductRepo(client)
	outboxRepo := &repo.OutboxRepo{}
	committer := committer.NewCommitter(client)
	readModel := repo.NewReadModel(client)

	createProduct := create_product.NewInteractor(productRepo, outboxRepo, committer, clk)
	updateProduct := update_product.NewInteractor(productRepo, outboxRepo, committer, clk)
	activateProduct := activate_product.NewInteractor(productRepo, outboxRepo, committer, clk)
	deactivateProduct := deactivate_product.NewInteractor(productRepo, outboxRepo, committer, clk)
	applyDiscount := apply_discount.NewInteractor(productRepo, outboxRepo, committer, clk)
	removeDiscount := remove_discount.NewInteractor(productRepo, outboxRepo, committer, clk)
	archiveProduct := archive_product.NewInteractor(productRepo, outboxRepo, committer, clk)

	getProductQuery := &get_product.QueryImpl{ReadModel: readModel}
	listProductsQuery := &list_products.QueryImpl{ReadModel: readModel}

	return &Options{
		SpannerClient:     client,
		CreateProduct:     createProduct,
		UpdateProduct:     updateProduct,
		ActivateProduct:   activateProduct,
		DeactivateProduct: deactivateProduct,
		ApplyDiscount:     applyDiscount,
		RemoveDiscount:    removeDiscount,
		ArchiveProduct:    archiveProduct,
		GetProductQuery:   getProductQuery,
		ListProductsQuery: listProductsQuery,
	}
}

// ProductHandler builds the gRPC handler with adapters.
func (o *Options) ProductHandler() *product.Handler {
	return &product.Handler{
		CreateProduct:     &createProductAdapter{o.CreateProduct},
		UpdateProduct:     &updateProductAdapter{o.UpdateProduct},
		ActivateProduct:   &activateProductAdapter{o.ActivateProduct},
		DeactivateProduct: &deactivateProductAdapter{o.DeactivateProduct},
		ApplyDiscount:     &applyDiscountAdapter{o.ApplyDiscount},
		RemoveDiscount:    &removeDiscountAdapter{o.RemoveDiscount},
		ArchiveProduct:    &archiveProductAdapter{o.ArchiveProduct},
		GetProduct:        &getProductAdapter{o.GetProductQuery},
		ListProducts:      &listProductsAdapter{o.ListProductsQuery},
	}
}

// Adapters that convert transport request types to usecase types and implement product.*Runner.

type createProductAdapter struct{ *create_product.Interactor }
func (a *createProductAdapter) Execute(ctx context.Context, req product.CreateProductRequest) (string, error) {
	return a.Interactor.Execute(ctx, create_product.Request{
		Name: req.Name, Description: req.Description, Category: req.Category,
		BasePriceNum: req.BasePriceNum, BasePriceDenom: req.BasePriceDenom,
	})
}

type updateProductAdapter struct{ *update_product.Interactor }
func (a *updateProductAdapter) Execute(ctx context.Context, req product.UpdateProductRequest) error {
	return a.Interactor.Execute(ctx, update_product.Request{
		ProductID: req.ProductID, Name: req.Name, Description: req.Description, Category: req.Category,
	})
}

type activateProductAdapter struct{ *activate_product.Interactor }
func (a *activateProductAdapter) Execute(ctx context.Context, req product.ActivateProductRequest) error {
	return a.Interactor.Execute(ctx, activate_product.Request{ProductID: req.ProductID})
}

type deactivateProductAdapter struct{ *deactivate_product.Interactor }
func (a *deactivateProductAdapter) Execute(ctx context.Context, req product.DeactivateProductRequest) error {
	return a.Interactor.Execute(ctx, deactivate_product.Request{ProductID: req.ProductID})
}

type applyDiscountAdapter struct{ *apply_discount.Interactor }
func (a *applyDiscountAdapter) Execute(ctx context.Context, req product.ApplyDiscountRequest) error {
	start := time.Unix(req.StartDateUnix, 0)
	end := time.Unix(req.EndDateUnix, 0)
	return a.Interactor.Execute(ctx, apply_discount.Request{
		ProductID: req.ProductID,
		Percent:   req.Percent,
		StartDate: start,
		EndDate:   end,
	})
}

type removeDiscountAdapter struct{ *remove_discount.Interactor }
func (a *removeDiscountAdapter) Execute(ctx context.Context, req product.RemoveDiscountRequest) error {
	return a.Interactor.Execute(ctx, remove_discount.Request{ProductID: req.ProductID})
}

type archiveProductAdapter struct{ *archive_product.Interactor }
func (a *archiveProductAdapter) Execute(ctx context.Context, req product.ArchiveProductRequest) error {
	return a.Interactor.Execute(ctx, archive_product.Request{ProductID: req.ProductID})
}

type getProductAdapter struct{ q get_product.Query }
func (a *getProductAdapter) Execute(ctx context.Context, productID string) (*product.GetProductDTO, error) {
	dto, err := a.q.Execute(ctx, productID)
	if err != nil {
		return nil, err
	}
	return &product.GetProductDTO{
		ProductID:       dto.ProductID,
		Name:            dto.Name,
		Description:     dto.Description,
		Category:        dto.Category,
		BasePrice:       product.RatToDecimalString(dto.BasePrice),
		EffectivePrice:  product.RatToDecimalString(dto.EffectivePrice),
		DiscountPercent: dto.DiscountPercent,
		Status:          dto.Status,
	}, nil
}

type listProductsAdapter struct{ list_products.Query }
func (a *listProductsAdapter) Execute(ctx context.Context, req product.ListProductsRequest) (*product.ListProductsResultDTO, error) {
	res, err := a.Query.Execute(ctx, list_products.ListProductsRequest{
		Category: req.Category, Status: req.Status, Limit: req.Limit, Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*product.ProductSummaryDTO, len(res.Products))
	for i, p := range res.Products {
		out[i] = &product.ProductSummaryDTO{
			ProductID:      p.ProductID,
			Name:           p.Name,
			Description:    p.Description,
			Category:       p.Category,
			BasePrice:      ratToStr(p.BasePrice),
			EffectivePrice: ratToStr(p.EffectivePrice),
			Status:         p.Status,
		}
	}
	return &product.ListProductsResultDTO{Products: out, Total: res.Total}, nil
}

func ratToStr(r *big.Rat) string {
	if r == nil {
		return "0"
	}
	f, _ := r.Float64()
	return strconv.FormatFloat(f, 'f', 2, 64)
}
