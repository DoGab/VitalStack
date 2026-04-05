package datasource_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dogab/vitalstack/api/pkg/datasource"
)

func TestFSVOClient_SearchProducts_Success(t *testing.T) {
	searchFixture := loadTestFixture(t, "testdata/fsvo_search.json")
	foodFixture := loadTestFixture(t, "testdata/fsvo_food.json")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/foods":
			if q := r.URL.Query().Get("search"); q != "broccoli" {
				t.Errorf("expected search=broccoli, got %s", q)
			}
			if lang := r.URL.Query().Get("lang"); lang != "de" {
				t.Errorf("expected lang=de, got %s", lang)
			}
			_, _ = w.Write(searchFixture)
		case strings.HasPrefix(r.URL.Path, "/food/"):
			_, _ = w.Write(foodFixture)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := datasource.NewFSVOClient(srv.Client(), srv.URL, datasource.WithFSVOLanguage("de"))

	products, err := client.SearchProducts(context.Background(), "broccoli", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}

	p := products[0]
	if p.ID != "fsvo-349687" {
		t.Errorf("expected ID fsvo-349687, got %s", p.ID)
	}
	if p.Name != "Broccoli, raw" {
		t.Errorf("expected name 'Broccoli, raw', got %s", p.Name)
	}
	if p.Source != "fsvo" {
		t.Errorf("expected source fsvo, got %s", p.Source)
	}
	if len(p.Categories) != 1 || p.Categories[0] != "Fresh vegetables" {
		t.Errorf("expected category 'Fresh vegetables', got %v", p.Categories)
	}
}

func TestFSVOClient_SearchProducts_NutrientMapping(t *testing.T) {
	searchFixture := loadTestFixture(t, "testdata/fsvo_search.json")
	foodFixture := loadTestFixture(t, "testdata/fsvo_food.json")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/foods" {
			_, _ = w.Write(searchFixture)
		} else {
			_, _ = w.Write(foodFixture)
		}
	}))
	defer srv.Close()

	client := datasource.NewFSVOClient(srv.Client(), srv.URL, datasource.WithFSVOLanguage("de"))

	products, err := client.SearchProducts(context.Background(), "broccoli", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	macros := products[0].Macros
	if macros.Calories != 39 {
		t.Errorf("expected 39 calories, got %f", macros.Calories)
	}
	if macros.Protein != 3.6 {
		t.Errorf("expected 3.6g protein, got %f", macros.Protein)
	}
	if macros.Carbs != 3.2 {
		t.Errorf("expected 3.2g carbs, got %f", macros.Carbs)
	}
	if macros.Fat != 0.6 {
		t.Errorf("expected 0.6g fat, got %f", macros.Fat)
	}
	if macros.Fiber != 3.2 {
		t.Errorf("expected 3.2g fiber, got %f", macros.Fiber)
	}
}

func TestFSVOClient_LookupBarcode_AlwaysNotFound(t *testing.T) {
	client := datasource.NewFSVOClient(http.DefaultClient, "http://unused", datasource.WithFSVOLanguage("de"))

	_, err := client.LookupBarcode(context.Background(), "0000000000000")
	if !errors.Is(err, datasource.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFSVOClient_SearchProducts_EmptyResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	client := datasource.NewFSVOClient(srv.Client(), srv.URL, datasource.WithFSVOLanguage("de"))

	products, err := client.SearchProducts(context.Background(), "zzzznonexistent", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestFSVOClient_Name(t *testing.T) {
	client := datasource.NewFSVOClient(http.DefaultClient, "http://unused")
	if client.Name() != "fsvo" {
		t.Errorf("expected name 'fsvo', got %s", client.Name())
	}
}
