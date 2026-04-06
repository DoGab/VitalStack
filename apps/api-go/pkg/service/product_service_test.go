package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dogab/vitalstack/api/pkg/datasource"
	"github.com/dogab/vitalstack/api/pkg/service"
	"github.com/dogab/vitalstack/api/pkg/types"
)

// mockSearchIndex is a test double for search.ProductSearchIndex.
type mockSearchIndex struct {
	lookupResult  *types.Product
	lookupErr     error
	searchResults []types.Product
	searchErr     error
	indexedSingle []types.Product
	indexedBulk   []types.Product
}

func (m *mockSearchIndex) IndexProduct(_ context.Context, product types.Product) error {
	m.indexedSingle = append(m.indexedSingle, product)
	return nil
}

func (m *mockSearchIndex) IndexProducts(_ context.Context, products []types.Product) error {
	m.indexedBulk = append(m.indexedBulk, products...)
	return nil
}

func (m *mockSearchIndex) Search(_ context.Context, _ string, _ int) ([]types.Product, error) {
	return m.searchResults, m.searchErr
}

func (m *mockSearchIndex) LookupBarcode(_ context.Context, _ string) (*types.Product, error) {
	return m.lookupResult, m.lookupErr
}

// mockDatasource is a test double for datasource.FoodDatasource.
type mockDatasource struct {
	name          string
	lookupResult  *types.Product
	lookupErr     error
	searchResults []types.Product
	searchErr     error
	lookupCalled  bool
	searchCalled  bool
}

func (m *mockDatasource) Name() string { return m.name }

func (m *mockDatasource) LookupBarcode(_ context.Context, _ string, _ string) (*types.Product, error) {
	m.lookupCalled = true
	return m.lookupResult, m.lookupErr
}

func (m *mockDatasource) SearchProducts(_ context.Context, _ string, _ int, _ string) ([]types.Product, error) {
	m.searchCalled = true
	return m.searchResults, m.searchErr
}

func TestProductService_LookupBarcode_CacheHit(t *testing.T) {
	cached := &types.Product{ID: "off-123", Barcode: "123", Name: "Cached Product", Source: "openfoodfacts"}
	idx := &mockSearchIndex{lookupResult: cached}
	ds := &mockDatasource{name: "off"}

	svc := service.NewProductService(idx, ds)

	product, err := svc.LookupBarcode(context.Background(), "123", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.ID != "off-123" {
		t.Errorf("expected cached product, got ID %s", product.ID)
	}
	if ds.lookupCalled {
		t.Error("external datasource should NOT be called on cache hit")
	}
}

func TestProductService_LookupBarcode_CacheMiss_FallsToExternal(t *testing.T) {
	idx := &mockSearchIndex{lookupResult: nil} // cache miss
	offProduct := &types.Product{ID: "off-456", Barcode: "456", Name: "OFF Product", Source: "openfoodfacts"}
	offDS := &mockDatasource{name: "off", lookupResult: offProduct}
	usdaDS := &mockDatasource{name: "usda", lookupErr: datasource.ErrNotFound}

	svc := service.NewProductService(idx, offDS, usdaDS)

	product, err := svc.LookupBarcode(context.Background(), "456", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.ID != "off-456" {
		t.Errorf("expected OFF product, got ID %s", product.ID)
	}
	if !offDS.lookupCalled {
		t.Error("OFF datasource should be called on cache miss")
	}
	if usdaDS.lookupCalled {
		t.Error("USDA should NOT be called when OFF succeeds")
	}

	// Wait briefly for async caching
	time.Sleep(50 * time.Millisecond)
	if len(idx.indexedSingle) != 1 {
		t.Errorf("expected 1 product indexed, got %d", len(idx.indexedSingle))
	}
}

func TestProductService_LookupBarcode_CascadesToUSDA(t *testing.T) {
	idx := &mockSearchIndex{lookupResult: nil}
	offDS := &mockDatasource{name: "off", lookupErr: datasource.ErrNotFound}
	usdaProduct := &types.Product{ID: "usda-789", Barcode: "789", Name: "USDA Product", Source: "usda"}
	usdaDS := &mockDatasource{name: "usda", lookupResult: usdaProduct}

	svc := service.NewProductService(idx, offDS, usdaDS)

	product, err := svc.LookupBarcode(context.Background(), "789", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.ID != "usda-789" {
		t.Errorf("expected USDA product, got ID %s", product.ID)
	}
	if !offDS.lookupCalled {
		t.Error("OFF should be tried first")
	}
	if !usdaDS.lookupCalled {
		t.Error("USDA should be tried after OFF miss")
	}
}

func TestProductService_LookupBarcode_AllMiss(t *testing.T) {
	idx := &mockSearchIndex{lookupResult: nil}
	offDS := &mockDatasource{name: "off", lookupErr: datasource.ErrNotFound}
	usdaDS := &mockDatasource{name: "usda", lookupErr: datasource.ErrNotFound}

	svc := service.NewProductService(idx, offDS, usdaDS)

	_, err := svc.LookupBarcode(context.Background(), "000", "")
	if !errors.Is(err, datasource.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestProductService_SearchProducts_IndexSufficient(t *testing.T) {
	indexProducts := []types.Product{
		{ID: "off-1", Barcode: "1", Name: "Product A"},
		{ID: "off-2", Barcode: "2", Name: "Product B"},
	}
	idx := &mockSearchIndex{searchResults: indexProducts}
	ds := &mockDatasource{name: "off"}

	svc := service.NewProductService(idx, ds)

	results, err := svc.SearchProducts(context.Background(), "product", 2, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if ds.searchCalled {
		t.Error("external datasource should NOT be called when index has enough results")
	}
}

func TestProductService_SearchProducts_FanOut(t *testing.T) {
	indexProducts := []types.Product{
		{ID: "off-1", Barcode: "1", Name: "Product A"},
	}
	idx := &mockSearchIndex{searchResults: indexProducts}

	offDS := &mockDatasource{
		name: "off",
		searchResults: []types.Product{
			{ID: "off-3", Barcode: "3", Name: "Product C"},
		},
	}
	usdaDS := &mockDatasource{
		name: "usda",
		searchResults: []types.Product{
			{ID: "usda-4", Barcode: "4", Name: "Product D"},
		},
	}

	svc := service.NewProductService(idx, offDS, usdaDS)

	results, err := svc.SearchProducts(context.Background(), "product", 5, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 merged results, got %d", len(results))
	}
}

func TestProductService_SearchProducts_Deduplication(t *testing.T) {
	indexProducts := []types.Product{
		{ID: "off-1", Barcode: "123", Name: "Product A"},
	}
	idx := &mockSearchIndex{searchResults: indexProducts}

	offDS := &mockDatasource{
		name: "off",
		searchResults: []types.Product{
			{ID: "off-1-dup", Barcode: "123", Name: "Product A (duplicate)"}, // same barcode
			{ID: "off-2", Barcode: "456", Name: "Product B"},
		},
	}

	svc := service.NewProductService(idx, offDS)

	results, err := svc.SearchProducts(context.Background(), "product", 10, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results (deduplicated), got %d", len(results))
	}
}
