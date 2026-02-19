package tmpl

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	expectedPages := []string{"home", "contacts", "contact-form"}
	for _, name := range expectedPages {
		if _, ok := r.pages[name]; !ok {
			t.Errorf("missing page template: %s", name)
		}
	}

	expectedPartials := []string{"contact-rows"}
	for _, name := range expectedPartials {
		if _, ok := r.partials[name]; !ok {
			t.Errorf("missing partial template: %s", name)
		}
	}
}

func TestRenderPage(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var buf bytes.Buffer
	err = r.RenderPage(&buf, "home", nil)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}

	body := buf.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("expected DOCTYPE in output")
	}
	if !strings.Contains(body, "htmxapp") {
		t.Error("expected htmxapp in output")
	}
}

func TestRenderPage_NotFound(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var buf bytes.Buffer
	err = r.RenderPage(&buf, "nonexistent", nil)
	if err == nil {
		t.Error("expected error for missing template")
	}
}

func TestRenderPartial(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	type contact struct {
		ID        string
		FirstName string
		LastName  string
		Email     string
		Phone     string
	}

	data := []contact{
		{ID: "1", FirstName: "Alice", LastName: "Smith", Email: "alice@test.com", Phone: "555-0001"},
	}

	var buf bytes.Buffer
	err = r.RenderPartial(&buf, "contact-rows", data)
	if err != nil {
		t.Fatalf("RenderPartial: %v", err)
	}

	body := buf.String()
	if !strings.Contains(body, "Alice") {
		t.Error("expected Alice in partial output")
	}
}

func TestExtractName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"templates/pages/home.html", "home"},
		{"templates/partials/contact-rows.html", "contact-rows"},
		{"layout.html", "layout"},
	}
	for _, tt := range tests {
		if got := extractName(tt.path); got != tt.want {
			t.Errorf("extractName(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}
