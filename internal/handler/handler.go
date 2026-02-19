package handler

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/devaloi/htmxapp/internal/store"
	"github.com/devaloi/htmxapp/internal/tmpl"
)

//go:embed static
var staticFS embed.FS

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	store    store.ContactStore
	renderer *tmpl.Renderer
}

// New creates a Handler with the given store and renderer.
func New(s store.ContactStore, r *tmpl.Renderer) *Handler {
	return &Handler{store: s, renderer: r}
}

// Routes returns an http.Handler with all routes registered.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	// Static files
	staticSub, _ := fs.Sub(staticFS, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	// Pages
	mux.HandleFunc("GET /{$}", h.Home)
	mux.HandleFunc("GET /contacts", h.ListContacts)
	mux.HandleFunc("GET /contacts/new", h.NewContact)
	mux.HandleFunc("POST /contacts", h.CreateContact)
	mux.HandleFunc("GET /contacts/search", h.SearchContacts)
	mux.HandleFunc("GET /contacts/{id}/edit", h.EditContact)
	mux.HandleFunc("POST /contacts/{id}", h.UpdateContact)
	mux.HandleFunc("DELETE /contacts/{id}", h.DeleteContact)

	return mux
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
