package model

import (
	"testing"
)

func TestContact_FullName(t *testing.T) {
	tests := []struct {
		name    string
		contact Contact
		want    string
	}{
		{"both names", Contact{FirstName: "Jane", LastName: "Doe"}, "Jane Doe"},
		{"first only", Contact{FirstName: "Jane"}, "Jane"},
		{"last only", Contact{LastName: "Doe"}, "Doe"},
		{"empty", Contact{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.contact.FullName(); got != tt.want {
				t.Errorf("FullName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestContact_Validate(t *testing.T) {
	tests := []struct {
		name      string
		contact   Contact
		wantErrs  []string
		wantClean bool
	}{
		{
			name:      "valid contact",
			contact:   Contact{FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"},
			wantClean: true,
		},
		{
			name:     "missing first name",
			contact:  Contact{LastName: "Doe", Email: "jane@example.com"},
			wantErrs: []string{"FirstName"},
		},
		{
			name:     "missing last name",
			contact:  Contact{FirstName: "Jane", Email: "jane@example.com"},
			wantErrs: []string{"LastName"},
		},
		{
			name:     "missing email",
			contact:  Contact{FirstName: "Jane", LastName: "Doe"},
			wantErrs: []string{"Email"},
		},
		{
			name:     "invalid email",
			contact:  Contact{FirstName: "Jane", LastName: "Doe", Email: "not-an-email"},
			wantErrs: []string{"Email"},
		},
		{
			name:     "all missing",
			contact:  Contact{},
			wantErrs: []string{"FirstName", "LastName", "Email"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.contact.Validate()
			if tt.wantClean {
				if len(errs) != 0 {
					t.Errorf("expected no errors, got %v", errs)
				}
				return
			}
			for _, key := range tt.wantErrs {
				if _, ok := errs[key]; !ok {
					t.Errorf("expected error for %q, got none", key)
				}
			}
		})
	}
}
