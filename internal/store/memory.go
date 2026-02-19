package store

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/devaloi/htmxapp/internal/model"
)

// Memory is a thread-safe in-memory contact store.
type Memory struct {
	mu      sync.RWMutex
	data    map[string]model.Contact
	emails  map[string]string // email -> id for uniqueness
	counter int
}

// NewMemory creates a new in-memory store.
func NewMemory() *Memory {
	return &Memory{
		data:   make(map[string]model.Contact),
		emails: make(map[string]string),
	}
}

func (m *Memory) nextID() string {
	m.counter++
	return fmt.Sprintf("%d", m.counter)
}

// List returns contacts matching the search query, sorted by last name.
func (m *Memory) List(_ context.Context, search string) ([]model.Contact, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(search))
	result := make([]model.Contact, 0, len(m.data))

	for _, c := range m.data {
		if q == "" || m.matches(c, q) {
			result = append(result, c)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].LastName == result[j].LastName {
			return result[i].FirstName < result[j].FirstName
		}
		return result[i].LastName < result[j].LastName
	})

	return result, nil
}

func (m *Memory) matches(c model.Contact, q string) bool {
	fields := []string{c.FirstName, c.LastName, c.Email, c.Phone}
	for _, f := range fields {
		if strings.Contains(strings.ToLower(f), q) {
			return true
		}
	}
	return false
}

// Get returns a contact by ID.
func (m *Memory) Get(_ context.Context, id string) (model.Contact, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.data[id]
	if !ok {
		return model.Contact{}, model.ErrNotFound
	}
	return c, nil
}

// Create adds a new contact to the store.
func (m *Memory) Create(_ context.Context, c model.Contact) (model.Contact, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	email := strings.ToLower(c.Email)
	if existingID, ok := m.emails[email]; ok {
		if existingID != c.ID {
			return model.Contact{}, model.ErrDuplicateEmail
		}
	}

	now := time.Now()
	c.ID = m.nextID()
	c.CreatedAt = now
	c.UpdatedAt = now

	m.data[c.ID] = c
	m.emails[email] = c.ID
	return c, nil
}

// Update modifies an existing contact.
func (m *Memory) Update(_ context.Context, c model.Contact) (model.Contact, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, ok := m.data[c.ID]
	if !ok {
		return model.Contact{}, model.ErrNotFound
	}

	email := strings.ToLower(c.Email)
	if existingID, ok := m.emails[email]; ok && existingID != c.ID {
		return model.Contact{}, model.ErrDuplicateEmail
	}

	// Remove old email mapping
	delete(m.emails, strings.ToLower(existing.Email))

	c.CreatedAt = existing.CreatedAt
	c.UpdatedAt = time.Now()

	m.data[c.ID] = c
	m.emails[email] = c.ID
	return c, nil
}

// Delete removes a contact by ID.
func (m *Memory) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.data[id]
	if !ok {
		return model.ErrNotFound
	}

	delete(m.emails, strings.ToLower(c.Email))
	delete(m.data, id)
	return nil
}

// Count returns the total number of contacts.
func (m *Memory) Count(_ context.Context) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Seed adds sample contacts for development.
func (m *Memory) Seed() {
	samples := []model.Contact{
		{FirstName: "Alice", LastName: "Johnson", Email: "alice@example.com", Phone: "555-0101"},
		{FirstName: "Bob", LastName: "Smith", Email: "bob@example.com", Phone: "555-0102"},
		{FirstName: "Carol", LastName: "Williams", Email: "carol@example.com", Phone: "555-0103"},
		{FirstName: "David", LastName: "Brown", Email: "david@example.com", Phone: "555-0104"},
		{FirstName: "Eve", LastName: "Davis", Email: "eve@example.com", Phone: "555-0105"},
	}
	ctx := context.Background()
	for _, c := range samples {
		_, _ = m.Create(ctx, c)
	}
}
