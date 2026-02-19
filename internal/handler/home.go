package handler

import "net/http"

// Home renders the landing page.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if err := h.renderer.RenderPage(w, "home", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
