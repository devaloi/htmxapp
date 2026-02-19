package store

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/devaloi/htmxapp/internal/model"
)

func newTestStore(t *testing.T) *Memory {
	t.Helper()
	s := NewMemory()
	ctx := context.Background()

	contacts := []model.Contact{
		{FirstName: "Alice", LastName: "Johnson", Email: "alice@example.com", Phone: "555-0001"},
		{FirstName: "Bob", LastName: "Smith", Email: "bob@example.com", Phone: "555-0002"},
		{FirstName: "Carol", LastName: "Williams", Email: "carol@example.com"},
	}
	for _, c := range contacts {
		_, err := s.Create(ctx, c)
		if err != nil {
			t.Fatalf("seed create: %v", err)
		}
	}
	return s
}

func TestMemory_Create(t *testing.T) {
	s := NewMemory()
	ctx := context.Background()

	c, err := s.Create(ctx, model.Contact{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@example.com",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if c.ID == "" {
		t.Error("expected non-empty ID")
	}
	if c.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestMemory_Create_DuplicateEmail(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.Create(ctx, model.Contact{
		FirstName: "Another",
		LastName:  "Alice",
		Email:     "alice@example.com",
	})
	if !errors.Is(err, model.ErrDuplicateEmail) {
		t.Errorf("expected ErrDuplicateEmail, got %v", err)
	}
}

func TestMemory_Get(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	c, err := s.Get(ctx, "1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if c.FirstName != "Alice" {
		t.Errorf("expected Alice, got %s", c.FirstName)
	}
}

func TestMemory_Get_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.Get(ctx, "999")
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemory_List(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	contacts, err := s.List(ctx, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(contacts) != 3 {
		t.Errorf("expected 3 contacts, got %d", len(contacts))
	}
	// Sorted by last name
	if contacts[0].LastName != "Johnson" {
		t.Errorf("expected Johnson first, got %s", contacts[0].LastName)
	}
}

func TestMemory_List_Search(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	contacts, err := s.List(ctx, "alice")
	if err != nil {
		t.Fatalf("List with search: %v", err)
	}
	if len(contacts) != 1 {
		t.Errorf("expected 1 result, got %d", len(contacts))
	}
}

func TestMemory_Update(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	updated, err := s.Update(ctx, model.Contact{
		ID:        "1",
		FirstName: "Alicia",
		LastName:  "Johnson",
		Email:     "alicia@example.com",
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.FirstName != "Alicia" {
		t.Errorf("expected Alicia, got %s", updated.FirstName)
	}
	if updated.CreatedAt.IsZero() {
		t.Error("expected CreatedAt preserved")
	}
}

func TestMemory_Update_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.Update(ctx, model.Contact{ID: "999", FirstName: "X", LastName: "Y", Email: "x@y.com"})
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemory_Update_DuplicateEmail(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.Update(ctx, model.Contact{
		ID:        "1",
		FirstName: "Alice",
		LastName:  "Johnson",
		Email:     "bob@example.com",
	})
	if !errors.Is(err, model.ErrDuplicateEmail) {
		t.Errorf("expected ErrDuplicateEmail, got %v", err)
	}
}

func TestMemory_Delete(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.Delete(ctx, "1")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = s.Get(ctx, "1")
	if !errors.Is(err, model.ErrNotFound) {
		t.Error("expected contact to be deleted")
	}

	if s.Count(ctx) != 2 {
		t.Errorf("expected 2 contacts, got %d", s.Count(ctx))
	}
}

func TestMemory_Delete_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.Delete(ctx, "999")
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemory_Count(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	if got := s.Count(ctx); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
}

func TestMemory_ConcurrentAccess(t *testing.T) {
	s := NewMemory()
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := range 50 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c := model.Contact{
				FirstName: "User",
				LastName:  "Test",
				Email:     fmt.Sprintf("user%d@example.com", n),
			}
			_, _ = s.Create(ctx, c)
		}(i)
	}

	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.List(ctx, "")
		}()
	}

	wg.Wait()

	if got := s.Count(ctx); got != 50 {
		t.Errorf("expected 50 contacts, got %d", got)
	}
}
