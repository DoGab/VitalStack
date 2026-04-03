package datasource_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dogab/vitalstack/api/pkg/datasource"
)

func loadTestFixture(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test fixture %s: %v", path, err)
	}
	return data
}

func TestOFFClient_LookupBarcode_Success(t *testing.T) {
	fixture := loadTestFixture(t, "testdata/off_product.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request path contains the barcode
		if r.URL.Path != "/api/v2/product/7613035466432" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// Verify User-Agent is set
		if ua := r.Header.Get("User-Agent"); ua == "" {
			t.Error("User-Agent header not set")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(fixture)
	}))
	defer srv.Close()

	client := datasource.NewOFFClient(srv.Client())
	client.SetBaseURL(srv.URL) // test helper

	product, err := client.LookupBarcode(context.Background(), "7613035466432")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.ID != "off-7613035466432" {
		t.Errorf("expected ID off-7613035466432, got %s", product.ID)
	}
	if product.Name != "Caffè Latte" {
		t.Errorf("expected name Caffè Latte, got %s", product.Name)
	}
	if product.Brand != "Emmi" {
		t.Errorf("expected brand Emmi, got %s", product.Brand)
	}
	if product.Source != "openfoodfacts" {
		t.Errorf("expected source openfoodfacts, got %s", product.Source)
	}
	if product.NutriScore != "c" {
		t.Errorf("expected nutri_score c, got %s", product.NutriScore)
	}
	if product.Macros.Calories != 56 {
		t.Errorf("expected calories 56, got %f", product.Macros.Calories)
	}
	if product.Macros.Protein != 3.2 {
		t.Errorf("expected protein 3.2, got %f", product.Macros.Protein)
	}
	if len(product.Categories) != 3 {
		t.Errorf("expected 3 categories, got %d", len(product.Categories))
	}
	// Verify "en:" prefix is stripped
	if len(product.Categories) > 0 && product.Categories[0] != "beverages" {
		t.Errorf("expected category 'beverages', got '%s'", product.Categories[0])
	}
}

func TestOFFClient_LookupBarcode_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status": 0, "product": {}}`))
	}))
	defer srv.Close()

	client := datasource.NewOFFClient(srv.Client())
	client.SetBaseURL(srv.URL)

	_, err := client.LookupBarcode(context.Background(), "0000000000000")
	if !errors.Is(err, datasource.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestOFFClient_LookupBarcode_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{malformed`))
	}))
	defer srv.Close()

	client := datasource.NewOFFClient(srv.Client())
	client.SetBaseURL(srv.URL)

	_, err := client.LookupBarcode(context.Background(), "123")
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestOFFClient_SearchProducts_Success(t *testing.T) {
	fixture := loadTestFixture(t, "testdata/off_search.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi/search.pl" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query().Get("search_terms")
		if q != "nutella" {
			t.Errorf("expected search_terms=nutella, got %s", q)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(fixture)
	}))
	defer srv.Close()

	client := datasource.NewOFFClient(srv.Client())
	client.SetBaseURL(srv.URL)

	products, err := client.SearchProducts(context.Background(), "nutella", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
	if products[0].Name != "Nutella" {
		t.Errorf("expected first product Nutella, got %s", products[0].Name)
	}
	if products[0].Macros.Calories != 539 {
		t.Errorf("expected 539 kcal, got %f", products[0].Macros.Calories)
	}
}

func TestOFFClient_SearchProducts_EmptyResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"products": []}`))
	}))
	defer srv.Close()

	client := datasource.NewOFFClient(srv.Client())
	client.SetBaseURL(srv.URL)

	products, err := client.SearchProducts(context.Background(), "zzzznonexistent", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestOFFClient_LookupBarcode_MissingNutriments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status": 1,
			"product": {
				"code": "123",
				"product_name": "Unknown Product",
				"brands": "",
				"categories_tags": [],
				"image_url": "",
				"nutriscore_grade": "",
				"nutriments": {}
			}
		}`))
	}))
	defer srv.Close()

	client := datasource.NewOFFClient(srv.Client())
	client.SetBaseURL(srv.URL)

	product, err := client.LookupBarcode(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// All macros should be zero when nutriments are missing
	if product.Macros.Calories != 0 || product.Macros.Protein != 0 {
		t.Errorf("expected zero macros for missing nutriments, got %+v", product.Macros)
	}
}
