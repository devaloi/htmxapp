package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/devaloi/htmxapp/internal/model"
)

type contactListData struct {
	Contacts []model.Contact
	Count    int
	Search   string
}

type contactFormData struct {
	Contact model.Contact
	Errors  map[string]string
}

// ListContacts renders the full contacts page.
func (h *Handler) ListContacts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	contacts, err := h.store.List(r.Context(), q)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := contactListData{
		Contacts: contacts,
		Count:    h.store.Count(r.Context()),
		Search:   q,
	}

	if err := h.renderer.RenderPage(w, "contacts", data); err != nil {
		slog.Error("render contacts page", "error", err)
	}
}

// SearchContacts returns a partial with matching contact rows (htmx).
func (h *Handler) SearchContacts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	contacts, err := h.store.List(r.Context(), q)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := h.renderer.RenderPartial(w, "contact-rows", contacts); err != nil {
		slog.Error("render contact rows", "error", err)
	}
}

// NewContact renders the new contact form.
func (h *Handler) NewContact(w http.ResponseWriter, r *http.Request) {
	data := contactFormData{Errors: make(map[string]string)}
	if err := h.renderer.RenderPage(w, "contact-form", data); err != nil {
		slog.Error("render new contact form", "error", err)
	}
}

// CreateContact handles the form submission for creating a contact.
func (h *Handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	c := contactFromForm(r)

	if errs := c.Validate(); len(errs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		data := contactFormData{Contact: c, Errors: errs}
		if err := h.renderer.RenderPage(w, "contact-form", data); err != nil {
			slog.Error("render form with errors", "error", err)
		}
		return
	}

	created, err := h.store.Create(r.Context(), c)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			data := contactFormData{
				Contact: c,
				Errors:  map[string]string{"Email": "A contact with this email already exists"},
			}
			if renderErr := h.renderer.RenderPage(w, "contact-form", data); renderErr != nil {
				slog.Error("render form with duplicate error", "error", renderErr)
			}
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("contact created", "id", created.ID, "name", created.FullName())
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

// EditContact renders the edit form for a contact.
func (h *Handler) EditContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := h.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := contactFormData{Contact: c, Errors: make(map[string]string)}
	if err := h.renderer.RenderPage(w, "contact-form", data); err != nil {
		slog.Error("render edit form", "error", err)
	}
}

// UpdateContact handles the form submission for updating a contact.
func (h *Handler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c := contactFromForm(r)
	c.ID = id

	if errs := c.Validate(); len(errs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		data := contactFormData{Contact: c, Errors: errs}
		if err := h.renderer.RenderPage(w, "contact-form", data); err != nil {
			slog.Error("render edit form with errors", "error", err)
		}
		return
	}

	updated, err := h.store.Update(r.Context(), c)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if errors.Is(err, model.ErrDuplicateEmail) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			data := contactFormData{
				Contact: c,
				Errors:  map[string]string{"Email": "A contact with this email already exists"},
			}
			if renderErr := h.renderer.RenderPage(w, "contact-form", data); renderErr != nil {
				slog.Error("render edit form with duplicate error", "error", renderErr)
			}
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("contact updated", "id", updated.ID, "name", updated.FullName())
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

// DeleteContact removes a contact and returns empty content for htmx swap.
func (h *Handler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.store.Delete(r.Context(), id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("contact deleted", "id", id)

	if isHTMX(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func contactFromForm(r *http.Request) model.Contact {
	return model.Contact{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Phone:     r.FormValue("phone"),
	}
}
