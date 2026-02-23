//go:build e2e

package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"product-catalog-service/internal/app/product/domain"
	"product-catalog-service/internal/app/product/domain/services"
	"product-catalog-service/internal/app/product/queries/get_product"
	"product-catalog-service/internal/app/product/queries/list_products"
	"product-catalog-service/internal/app/product/repo"
	"product-catalog-service/internal/app/product/usecases/activate_product"
	"product-catalog-service/internal/app/product/usecases/apply_discount"
	"product-catalog-service/internal/app/product/usecases/create_product"
	"product-catalog-service/internal/app/product/usecases/update_product"
	"product-catalog-service/internal/pkg/clock"
	"product-catalog-service/internal/pkg/committer"
)

const (
	testDBPath = "projects/test-project/instances/test-instance/databases/product-catalog"
)

func TestProductCreationFlow(t *testing.T) {
	client, cleanup := setupSpanner(t)
	defer cleanup()

	clk := clock.RealClock{}
	productRepo := repo.NewProductRepo(client)
	outboxRepo := &repo.OutboxRepo{}
	comm := committer.NewCommitter(client)
	readModel := repo.NewReadModel(client)

	createUsecase := create_product.NewInteractor(productRepo, outboxRepo, comm, clk)
	getQuery := &get_product.QueryImpl{ReadModel: readModel}

	ctx := context.Background()

	// Create product
	productID, err := createUsecase.Execute(ctx, create_product.Request{
		Name:             "Test Product",
		Description:      "A test product",
		Category:         "electronics",
		BasePriceNum:     1999,
		BasePriceDenom:   100,
	})
	require.NoError(t, err)
	require.NotEmpty(t, productID)

	// Verify query returns correct data
	product, err := getQuery.Execute(ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, "electronics", product.Category)
	assert.Equal(t, "19.99", product.BasePrice.FloatString(2))

	// Verify outbox event was created
	events := getOutboxEvents(t, client, productID)
	require.Len(t, events, 1)
	assert.Equal(t, "product.created", events[0].EventType)
}

func TestDiscountApplicationFlow(t *testing.T) {
	client, cleanup := setupSpanner(t)
	defer cleanup()

	clk := clock.RealClock{}
	productRepo := repo.NewProductRepo(client)
	outboxRepo := &repo.OutboxRepo{}
	comm := committer.NewCommitter(client)
	readModel := repo.NewReadModel(client)

	createUsecase := create_product.NewInteractor(productRepo, outboxRepo, comm, clk)
	activateUsecase := activate_product.NewInteractor(productRepo, outboxRepo, comm, clk)
	applyDiscountUsecase := apply_discount.NewInteractor(productRepo, outboxRepo, comm, clk)
	getQuery := &get_product.QueryImpl{ReadModel: readModel}

	ctx := context.Background()

	productID, err := createUsecase.Execute(ctx, create_product.Request{
		Name:             "Discounted Product",
		Description:      "Product with discount",
		Category:         "books",
		BasePriceNum:     5000, // 50.00
		BasePriceDenom:   100,
	})
	require.NoError(t, err)
	require.NotEmpty(t, productID)

	err = activateUsecase.Execute(ctx, activate_product.Request{ProductID: productID})
	require.NoError(t, err)

	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now.Add(7 * 24 * time.Hour)

	err = applyDiscountUsecase.Execute(ctx, apply_discount.Request{
		ProductID: productID,
		Percent:   20,
		StartDate: start,
		EndDate:   end,
	})
	require.NoError(t, err)

	product, err := getQuery.Execute(ctx, productID)
	require.NoError(t, err)
	// 20% off 50.00 = 40.00
	expected := services.PricingCalculator_EffectivePrice(
		domain.NewMoney(5000, 100),
		domain.NewDiscount(20, start, end),
		domain.PercentToRat(20),
	)
	require.NotNil(t, expected)
	assert.Equal(t, expected.Rat().FloatString(2), product.EffectivePrice.FloatString(2))
}

func TestBusinessRuleValidation(t *testing.T) {
	client, cleanup := setupSpanner(t)
	defer cleanup()

	clk := clock.RealClock{}
	productRepo := repo.NewProductRepo(client)
	outboxRepo := &repo.OutboxRepo{}
	comm := committer.NewCommitter(client)

	createUsecase := create_product.NewInteractor(productRepo, outboxRepo, comm, clk)
	applyDiscountUsecase := apply_discount.NewInteractor(productRepo, outboxRepo, comm, clk)

	ctx := context.Background()

	productID, err := createUsecase.Execute(ctx, create_product.Request{
		Name:             "Draft Product",
		Category:         "test",
		BasePriceNum:     1000,
		BasePriceDenom:   100,
	})
	require.NoError(t, err)

	// Cannot apply discount to draft (inactive) product
	now := time.Now()
	err = applyDiscountUsecase.Execute(ctx, apply_discount.Request{
		ProductID: productID,
		Percent:   10,
		StartDate: now,
		EndDate:   now.Add(24 * time.Hour),
	})
	assert.ErrorIs(t, err, domain.ErrProductNotActive)
}

func TestProductUpdateFlow(t *testing.T) {
	client, cleanup := setupSpanner(t)
	defer cleanup()

	clk := clock.RealClock{}
	productRepo := repo.NewProductRepo(client)
	outboxRepo := &repo.OutboxRepo{}
	comm := committer.NewCommitter(client)
	readModel := repo.NewReadModel(client)

	createUsecase := create_product.NewInteractor(productRepo, outboxRepo, comm, clk)
	updateUsecase := update_product.NewInteractor(productRepo, outboxRepo, comm, clk)
	getQuery := &get_product.QueryImpl{ReadModel: readModel}

	ctx := context.Background()

	productID, err := createUsecase.Execute(ctx, create_product.Request{
		Name:             "Original Name",
		Description:      "Original desc",
		Category:         "cat",
		BasePriceNum:     100,
		BasePriceDenom:   1,
	})
	require.NoError(t, err)

	err = updateUsecase.Execute(ctx, update_product.Request{
		ProductID:   productID,
		Name:        "Updated Name",
		Description: "Updated desc",
		Category:    "cat",
	})
	require.NoError(t, err)

	product, err := getQuery.Execute(ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", product.Name)
	assert.Equal(t, "Updated desc", product.Description)
}

func TestListProducts(t *testing.T) {
	client, cleanup := setupSpanner(t)
	defer cleanup()

	readModel := repo.NewReadModel(client)
	listQuery := &list_products.QueryImpl{ReadModel: readModel}
	ctx := context.Background()

	res, err := listQuery.Execute(ctx, list_products.ListProductsRequest{
		Status: "active",
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Products)
}

func setupSpanner(t *testing.T) (*spanner.Client, func()) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("Set SPANNER_EMULATOR_HOST to run E2E tests (e.g. localhost:9010)")
	}
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, testDBPath)
	require.NoError(t, err)
	return client, func() { client.Close() }
}

func getOutboxEvents(t *testing.T, client *spanner.Client, aggregateID string) []struct {
	EventType string
} {
	ctx := context.Background()
	iter := client.Single().Query(ctx, spanner.Statement{
		SQL:    "SELECT event_type FROM outbox_events WHERE aggregate_id = @aid",
		Params: map[string]interface{}{"aid": aggregateID},
	})
	defer iter.Stop()

	var events []struct {
		EventType string
	}
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		require.NoError(t, err)
		var e struct {
			EventType string
		}
		err = row.Columns(&e.EventType)
		require.NoError(t, err)
		events = append(events, e)
	}
	return events
}
