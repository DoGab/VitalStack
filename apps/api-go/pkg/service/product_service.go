package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/dogab/vitalstack/api/pkg/datasource"
	"github.com/dogab/vitalstack/api/pkg/search"
	"github.com/dogab/vitalstack/api/pkg/types"
)

// ProductService orchestrates product lookup using a waterfall strategy:
// 1. Check Meilisearch index (cache)
// 2. Try external datasources in order (OFF → USDA)
// 3. Cache hits from external sources for future lookups
type ProductService struct {
	index       search.ProductSearchIndex
	datasources []datasource.FoodDatasource
}

// NewProductService creates a new ProductService with the given search index and datasources.
// Datasources are tried in order (waterfall pattern).
func NewProductService(index search.ProductSearchIndex, datasources ...datasource.FoodDatasource) *ProductService {
	return &ProductService{
		index:       index,
		datasources: datasources,
	}
}

// LookupBarcode searches all configured datasources in waterfall order.
// It returns the first matching product, or ErrNotFound if no datasource has a match.
// The lang parameter is forwarded to each datasource for localized results.
func (s *ProductService) LookupBarcode(ctx context.Context, barcode string, lang string) (*types.Product, error) {
	// Step 1: Check the search index (cache)
	product, err := s.index.LookupBarcode(ctx, barcode)
	if err != nil {
		slog.Warn("ProductService: index barcode lookup failed", "barcode", barcode, "error", err)
	}
	if product != nil {
		slog.Debug("ProductService: barcode cache hit", "barcode", barcode, "source", product.Source)
		return product, nil
	}

	// Step 2: Waterfall through external datasources
	for _, ds := range s.datasources {
		product, err = ds.LookupBarcode(ctx, barcode, lang)
		if err != nil {
			if errors.Is(err, datasource.ErrNotFound) {
				slog.Debug("ProductService: barcode not found in datasource", "barcode", barcode, "datasource", ds.Name())
				continue
			}
			slog.Warn("ProductService: datasource barcode lookup failed", "barcode", barcode, "datasource", ds.Name(), "error", err)
			continue
		}

		// Step 3: Cache the product in Meilisearch asynchronously
		s.cacheProductAsync(product)

		return product, nil
	}

	return nil, datasource.ErrNotFound
}

// SearchProducts performs a product search using the waterfall strategy.
func (s *ProductService) SearchProducts(ctx context.Context, query string, limit int, lang string) ([]types.Product, error) {
	// Step 1: Search the index
	results, err := s.index.Search(ctx, query, limit)
	if err != nil {
		slog.Warn("ProductService: index search failed", "query", query, "error", err)
		results = nil
	}

	// Step 2: If we have enough results from the index, return them
	if len(results) >= limit {
		return results[:limit], nil
	}

	// Step 3: Fan out to external datasources for more results
	remaining := limit - len(results)
	seenBarcodes := make(map[string]bool, len(results))
	for _, p := range results {
		if p.Barcode != "" {
			seenBarcodes[p.Barcode] = true
		}
	}

	var newProducts []types.Product
	for _, ds := range s.datasources {
		external, err := ds.SearchProducts(ctx, query, remaining, lang)
		if err != nil {
			slog.Warn("ProductService: external search failed", "query", query, "datasource", ds.Name(), "error", err)
			continue
		}

		// Deduplicate by barcode
		for _, p := range external {
			if p.Barcode != "" && seenBarcodes[p.Barcode] {
				continue
			}
			if p.Barcode != "" {
				seenBarcodes[p.Barcode] = true
			}
			results = append(results, p)
			newProducts = append(newProducts, p)
		}
	}

	// Step 4: Cache new products asynchronously
	if len(newProducts) > 0 {
		s.cacheProductsAsync(newProducts)
	}

	// Cap results
	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// cacheProductAsync asynchronously indexes a single product.
func (s *ProductService) cacheProductAsync(product *types.Product) {
	go func() {
		if err := s.index.IndexProduct(context.Background(), *product); err != nil {
			slog.Warn("ProductService: failed to cache product", "id", product.ID, "error", err)
		}
	}()
}

// cacheProductsAsync asynchronously indexes multiple products.
func (s *ProductService) cacheProductsAsync(products []types.Product) {
	// Copy slice to prevent data races
	productsCopy := make([]types.Product, len(products))
	copy(productsCopy, products)

	go func() {
		// Use a WaitGroup in case we want to extend this to batch
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.index.IndexProducts(context.Background(), productsCopy); err != nil {
				slog.Warn("ProductService: failed to cache products", "count", len(productsCopy), "error", err)
			}
		}()
		wg.Wait()
	}()
}
