package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AdamBrutsaert/basic-go-http-server/internal/store"
)

func TestItemHandlers(t *testing.T) {
	store := store.New()
	router := newRouter(store)

	req := httptest.NewRequest("POST", "/items", strings.NewReader(`{"name":"item1","price":100}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 201 {
		t.Errorf("expected status code 201, got %d", rec.Code)
	}

	req = httptest.NewRequest("GET", "/items", nil)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected status code 200, got %d", rec.Code)
	}

	expected := `[{"name":"item1","price":100}]` + "\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, rec.Body.String())
	}

	req = httptest.NewRequest("GET", "/items/1", nil)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected status code 200, got %d", rec.Code)
	}

	expected = `{"name":"item1","price":100}` + "\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, rec.Body.String())
	}

	req = httptest.NewRequest("PUT", "/items/1", strings.NewReader(`{"name":"item2","price":200}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 204 {
		t.Errorf("expected status code 204, got %d", rec.Code)
	}

	req = httptest.NewRequest("GET", "/items/1", nil)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected status code 200, got %d", rec.Code)
	}

	expected = `{"name":"item2","price":200}` + "\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, rec.Body.String())
	}

	req = httptest.NewRequest("DELETE", "/items/1", nil)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 204 {
		t.Errorf("expected status code 204, got %d", rec.Code)
	}

	req = httptest.NewRequest("GET", "/items/1", nil)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Errorf("expected status code 404, got %d", rec.Code)
	}

	expected = "item not found\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, rec.Body.String())
	}

	req = httptest.NewRequest("GET", "/items", nil)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected status code 200, got %d", rec.Code)
	}

	expected = "[]\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, rec.Body.String())
	}
}
