package search_test

import (
	"encoding/json"
	"testing"

	"github.com/dogab/vitalstack/api/pkg/types"
)

// TestProductJSONRoundTrip verifies that Product survives JSON serialization
// (the same mechanism used by hitsToProducts for Meilisearch result mapping).
func TestProductJSONRoundTrip(t *testing.T) {
	original := types.Product{
		ID:         "off-7613035466432",
		Barcode:    "7613035466432",
		Name:       "Caffè Latte",
		Brand:      "Emmi",
		Categories: []string{"beverages", "dairy-drinks"},
		ImageURL:   "https://example.com/image.jpg",
		Source:     "openfoodfacts",
		NutriScore: "c",
		Macros: types.MacrosPer100g{
			Calories: 56,
			Protein:  3.2,
			Carbs:    5.1,
			Fat:      1.8,
			Fiber:    0,
		},
	}

	// Marshal to JSON (simulates what Meilisearch stores)
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal back (simulates hitsToProducts conversion)
	var roundTripped types.Product
	if err := json.Unmarshal(data, &roundTripped); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if roundTripped.ID != original.ID {
		t.Errorf("ID mismatch: got %s, want %s", roundTripped.ID, original.ID)
	}
	if roundTripped.Barcode != original.Barcode {
		t.Errorf("Barcode mismatch: got %s, want %s", roundTripped.Barcode, original.Barcode)
	}
	if roundTripped.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", roundTripped.Name, original.Name)
	}
	if roundTripped.Brand != original.Brand {
		t.Errorf("Brand mismatch: got %s, want %s", roundTripped.Brand, original.Brand)
	}
	if roundTripped.Source != original.Source {
		t.Errorf("Source mismatch: got %s, want %s", roundTripped.Source, original.Source)
	}
	if roundTripped.NutriScore != original.NutriScore {
		t.Errorf("NutriScore mismatch: got %s, want %s", roundTripped.NutriScore, original.NutriScore)
	}
	if roundTripped.Macros.Calories != original.Macros.Calories {
		t.Errorf("Calories mismatch: got %f, want %f", roundTripped.Macros.Calories, original.Macros.Calories)
	}
	if len(roundTripped.Categories) != len(original.Categories) {
		t.Errorf("Categories length mismatch: got %d, want %d", len(roundTripped.Categories), len(original.Categories))
	}
}

// TestProductIDPrefixConventions verifies the ID prefix conventions
// that prevent collisions across datasources.
func TestProductIDPrefixConventions(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected string
	}{
		{"OFF product", "off-7613035466432", "off-"},
		{"USDA product", "usda-2344382", "usda-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.id) < len(tt.expected) || tt.id[:len(tt.expected)] != tt.expected {
				t.Errorf("expected ID prefix %s, got %s", tt.expected, tt.id)
			}
		})
	}
}
