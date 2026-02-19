package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/devaloi/htmxapp/internal/store"
	"github.com/devaloi/htmxapp/internal/tmpl"
)

func setupTestHandler(t *testing.T) (*Handler, *store.Memory) {
	t.Helper()
	renderer, err := tmpl.New()
	if err != nil {
		t.Fatalf("tmpl.New: %v", err)
	}
	s := store.NewMemory()
	s.Seed()
	return New(s, renderer), s
}

func TestHome(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "htmxapp") {
		t.Error("expected 'htmxapp' in home page")
	}
}

func TestListContacts(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Alice") {
		t.Error("expected seeded contacts in list")
	}
}

func TestSearchContacts(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodGet, "/contacts/search?q=alice", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Alice") {
		t.Error("expected Alice in search results")
	}
	if strings.Contains(body, "Bob") {
		t.Error("did not expect Bob in search for 'alice'")
	}
}

func TestNewContact(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodGet, "/contacts/new", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "New Contact") {
		t.Error("expected form page")
	}
}

func TestCreateContact(t *testing.T) {
	h, s := setupTestHandler(t)
	mux := h.Routes()

	form := url.Values{
		"first_name": {"Frank"},
		"last_name":  {"Test"},
		"email":      {"frank@example.com"},
		"phone":      {"555-9999"},
	}
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/contacts" {
		t.Errorf("expected redirect to /contacts, got %s", loc)
	}

	ctx := req.Context()
	if s.Count(ctx) != 6 {
		t.Errorf("expected 6 contacts after create, got %d", s.Count(ctx))
	}
}

func TestCreateContact_ValidationError(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	form := url.Values{
		"first_name": {""},
		"last_name":  {""},
		"email":      {""},
	}
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rec.Code)
	}
}

func TestCreateContact_DuplicateEmail(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	form := url.Values{
		"first_name": {"Dupe"},
		"last_name":  {"User"},
		"email":      {"alice@example.com"},
	}
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rec.Code)
	}
}

func TestEditContact(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodGet, "/contacts/1/edit", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Edit Contact") {
		t.Error("expected edit form")
	}
}

func TestEditContact_NotFound(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodGet, "/contacts/999/edit", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestUpdateContact(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	form := url.Values{
		"first_name": {"Alicia"},
		"last_name":  {"Johnson"},
		"email":      {"alicia@example.com"},
		"phone":      {"555-0001"},
		"_method":    {"PUT"},
	}
	req := httptest.NewRequest(http.MethodPost, "/contacts/1", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}

func TestUpdateContact_NotFound(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	form := url.Values{
		"first_name": {"X"},
		"last_name":  {"Y"},
		"email":      {"x@y.com"},
		"_method":    {"PUT"},
	}
	req := httptest.NewRequest(http.MethodPost, "/contacts/999", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestDeleteContact_HTMX(t *testing.T) {
	h, s := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodDelete, "/contacts/1", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	ctx := req.Context()
	if s.Count(ctx) != 4 {
		t.Errorf("expected 4 contacts after delete, got %d", s.Count(ctx))
	}
}

func TestDeleteContact_NotFound(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodDelete, "/contacts/999", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestDeleteContact_Redirect(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodDelete, "/contacts/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
}

func TestStaticFiles(t *testing.T) {
	h, _ := setupTestHandler(t)
	mux := h.Routes()

	req := httptest.NewRequest(http.MethodGet, "/static/css/style.css", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "css") {
		t.Errorf("expected css content type, got %s", ct)
	}
}
