// Package datasource defines the contract and implementations for external food data providers.
// Each datasource normalizes external API responses into the common types.Product model.
package datasource

import (
	"context"
	"errors"

	"github.com/dogab/vitalstack/api/pkg/types"
)

// ErrNotFound is returned when a barcode lookup finds no matching product.
// This is a sentinel error that allows the waterfall pattern to continue
// to the next datasource without logging a spurious error.
var ErrNotFound = errors.New("product not found")

// FoodDatasource defines the contract for any external food data provider.
type FoodDatasource interface {
	// LookupBarcode searches for a product by its EAN/UPC barcode.
	// Returns ErrNotFound if no product matches the barcode.
	LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)

	// SearchProducts performs a free-text search and returns up to limit results.
	SearchProducts(ctx context.Context, query string, limit int) ([]types.Product, error)

	// Name returns the datasource identifier (e.g. "openfoodfacts", "usda").
	Name() string
}
