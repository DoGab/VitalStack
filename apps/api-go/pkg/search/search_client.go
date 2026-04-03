// Package search provides the product search index abstraction and Meilisearch implementation.
// It indexes and queries normalized Product data for fast, typo-tolerant product lookup.
package search

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"

	"github.com/dogab/vitalstack/api/pkg/types"
)

const (
	productsIndexName = "products"
	primaryKeyField   = "id"
)

// ProductSearchIndex defines search and indexing operations for the product catalog.
type ProductSearchIndex interface {
	// IndexProduct upserts a single product into the search index.
	IndexProduct(ctx context.Context, product types.Product) error

	// IndexProducts upserts multiple products into the search index.
	IndexProducts(ctx context.Context, products []types.Product) error

	// Search performs a full-text search and returns up to limit results.
	Search(ctx context.Context, query string, limit int) ([]types.Product, error)

	// LookupBarcode searches the index for a product with the given barcode.
	LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)
}

// MeilisearchClient implements ProductSearchIndex using Meilisearch.
type MeilisearchClient struct {
	client meilisearch.ServiceManager
	index  meilisearch.IndexManager
}

// NewMeilisearchClient creates a new Meilisearch client and configures the products index.
// It sets searchable, filterable, and sortable attributes for optimal product search.
func NewMeilisearchClient(url, apiKey string) (*MeilisearchClient, error) {
	client := meilisearch.New(url, meilisearch.WithAPIKey(apiKey))

	// Create or get the products index
	_, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        productsIndexName,
		PrimaryKey: primaryKeyField,
	})
	if err != nil {
		// Index may already exist — that's OK
		_ = err
	}

	index := client.Index(productsIndexName)

	// Configure searchable attributes
	_, err = index.UpdateSearchableAttributes(&[]string{"name", "brand", "categories"})
	if err != nil {
		return nil, fmt.Errorf("meilisearch: configuring searchable attributes: %w", err)
	}

	// Configure filterable attributes (for barcode lookup and source filtering)
	filterAttrs := []any{"barcode", "source"}
	_, err = index.UpdateFilterableAttributes(&filterAttrs)
	if err != nil {
		return nil, fmt.Errorf("meilisearch: configuring filterable attributes: %w", err)
	}

	// Configure sortable attributes
	_, err = index.UpdateSortableAttributes(&[]string{"name"})
	if err != nil {
		return nil, fmt.Errorf("meilisearch: configuring sortable attributes: %w", err)
	}

	return &MeilisearchClient{
		client: client,
		index:  index,
	}, nil
}

// IndexProduct upserts a single product into the search index.
func (m *MeilisearchClient) IndexProduct(_ context.Context, product types.Product) error {
	_, err := m.index.AddDocuments([]types.Product{product}, nil)
	if err != nil {
		return fmt.Errorf("meilisearch: indexing product: %w", err)
	}
	return nil
}

// IndexProducts upserts multiple products into the search index.
func (m *MeilisearchClient) IndexProducts(_ context.Context, products []types.Product) error {
	if len(products) == 0 {
		return nil
	}
	_, err := m.index.AddDocuments(products, nil)
	if err != nil {
		return fmt.Errorf("meilisearch: bulk indexing products: %w", err)
	}
	return nil
}

// Search performs a full-text search and returns up to limit results.
func (m *MeilisearchClient) Search(_ context.Context, query string, limit int) ([]types.Product, error) {
	result, err := m.index.Search(query, &meilisearch.SearchRequest{
		Limit: int64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("meilisearch: searching: %w", err)
	}

	return hitsToProducts(result.Hits)
}

// LookupBarcode searches the index for a product with the given barcode.
func (m *MeilisearchClient) LookupBarcode(_ context.Context, barcode string) (*types.Product, error) {
	result, err := m.index.Search("", &meilisearch.SearchRequest{
		Filter: fmt.Sprintf("barcode = '%s'", barcode),
		Limit:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("meilisearch: barcode lookup: %w", err)
	}

	products, err := hitsToProducts(result.Hits)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil // not found in cache
	}

	return &products[0], nil
}

// hitsToProducts converts Meilisearch search hits to typed Product values.
// Hits are returned as any maps, so we round-trip through JSON.
func hitsToProducts(hits meilisearch.Hits) ([]types.Product, error) {
	if len(hits) == 0 {
		return nil, nil
	}

	products := make([]types.Product, 0, len(hits))
	for _, hit := range hits {
		data, err := json.Marshal(hit)
		if err != nil {
			return nil, fmt.Errorf("meilisearch: marshaling hit: %w", err)
		}

		var p types.Product
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("meilisearch: unmarshaling hit: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}
