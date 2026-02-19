package model

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// Contact represents a person in the contact list.
type Contact struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FullName returns the contact's full name.
func (c Contact) FullName() string {
	return strings.TrimSpace(c.FirstName + " " + c.LastName)
}

// Validate checks required fields and returns any validation errors.
func (c Contact) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(c.FirstName) == "" {
		errs["FirstName"] = "First name is required"
	}
	if strings.TrimSpace(c.LastName) == "" {
		errs["LastName"] = "Last name is required"
	}
	if strings.TrimSpace(c.Email) == "" {
		errs["Email"] = "Email is required"
	} else if _, err := mail.ParseAddress(c.Email); err != nil {
		errs["Email"] = fmt.Sprintf("Invalid email address: %s", c.Email)
	}
	return errs
}
