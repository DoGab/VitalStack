package datasource_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dogab/vitalstack/api/pkg/datasource"
)

func TestUSDAClient_SearchProducts_Success(t *testing.T) {
	fixture := loadTestFixture(t, "testdata/usda_search.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/foods/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(fixture)
	}))
	defer srv.Close()

	client := datasource.NewUSDAClient(srv.Client(), "test-api-key")
	client.SetBaseURL(srv.URL)

	products, err := client.SearchProducts(context.Background(), "yogurt", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}

	p := products[0]
	if p.ID != "usda-2344382" {
		t.Errorf("expected ID usda-2344382, got %s", p.ID)
	}
	if p.Name != "Greek Yogurt, Vanilla" {
		t.Errorf("expected name 'Greek Yogurt, Vanilla', got %s", p.Name)
	}
	if p.Brand != "Chobani" {
		t.Errorf("expected brand Chobani, got %s", p.Brand)
	}
	if p.Source != "usda" {
		t.Errorf("expected source usda, got %s", p.Source)
	}
}

func TestUSDAClient_SearchProducts_NutrientMapping(t *testing.T) {
	fixture := loadTestFixture(t, "testdata/usda_search.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(fixture)
	}))
	defer srv.Close()

	client := datasource.NewUSDAClient(srv.Client(), "test-key")
	client.SetBaseURL(srv.URL)

	products, err := client.SearchProducts(context.Background(), "yogurt", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	macros := products[0].Macros
	if macros.Calories != 97 {
		t.Errorf("expected 97 calories, got %f", macros.Calories)
	}
	if macros.Protein != 12.5 {
		t.Errorf("expected 12.5g protein, got %f", macros.Protein)
	}
	if macros.Carbs != 12 {
		t.Errorf("expected 12g carbs, got %f", macros.Carbs)
	}
	if macros.Fat != 0 {
		t.Errorf("expected 0g fat, got %f", macros.Fat)
	}
	if macros.Fiber != 0 {
		t.Errorf("expected 0g fiber, got %f", macros.Fiber)
	}
}

func TestUSDAClient_APIKeyInjection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.URL.Query().Get("api_key")
		if apiKey != "my-secret-key" {
			t.Errorf("expected api_key=my-secret-key, got %s", apiKey)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"foods": []}`))
	}))
	defer srv.Close()

	client := datasource.NewUSDAClient(srv.Client(), "my-secret-key")
	client.SetBaseURL(srv.URL)

	_, err := client.SearchProducts(context.Background(), "test", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUSDAClient_LookupBarcode_Success(t *testing.T) {
	fixture := loadTestFixture(t, "testdata/usda_search.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify dataType filter
		dataType := r.URL.Query().Get("dataType")
		if dataType != "Branded" {
			t.Errorf("expected dataType=Branded, got %s", dataType)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(fixture)
	}))
	defer srv.Close()

	client := datasource.NewUSDAClient(srv.Client(), "test-key")
	client.SetBaseURL(srv.URL)

	product, err := client.LookupBarcode(context.Background(), "0894700010045")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.Barcode != "0894700010045" {
		t.Errorf("expected barcode 0894700010045, got %s", product.Barcode)
	}
}

func TestUSDAClient_LookupBarcode_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"foods": []}`))
	}))
	defer srv.Close()

	client := datasource.NewUSDAClient(srv.Client(), "test-key")
	client.SetBaseURL(srv.URL)

	_, err := client.LookupBarcode(context.Background(), "nonexistent")
	if !errors.Is(err, datasource.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUSDAClient_SearchProducts_EmptyResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"foods": []}`))
	}))
	defer srv.Close()

	client := datasource.NewUSDAClient(srv.Client(), "test-key")
	client.SetBaseURL(srv.URL)

	products, err := client.SearchProducts(context.Background(), "zzznonexistent", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}
