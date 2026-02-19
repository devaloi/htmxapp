package tmpl

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
)

//go:embed templates
var templateFS embed.FS

// Renderer loads and renders HTML templates.
type Renderer struct {
	pages    map[string]*template.Template
	partials map[string]*template.Template
}

// New parses all templates from the embedded filesystem.
func New() (*Renderer, error) {
	r := &Renderer{
		pages:    make(map[string]*template.Template),
		partials: make(map[string]*template.Template),
	}

	layoutContent, err := fs.ReadFile(templateFS, "templates/layout.html")
	if err != nil {
		return nil, fmt.Errorf("reading layout: %w", err)
	}

	// Parse page templates (each extends layout)
	pageFiles, err := fs.Glob(templateFS, "templates/pages/*.html")
	if err != nil {
		return nil, fmt.Errorf("globbing pages: %w", err)
	}
	for _, path := range pageFiles {
		name := extractName(path)
		t, parseErr := template.New("layout").Parse(string(layoutContent))
		if parseErr != nil {
			return nil, fmt.Errorf("parsing layout for %s: %w", name, parseErr)
		}
		pageContent, readErr := fs.ReadFile(templateFS, path)
		if readErr != nil {
			return nil, fmt.Errorf("reading %s: %w", path, readErr)
		}
		if _, parseErr = t.Parse(string(pageContent)); parseErr != nil {
			return nil, fmt.Errorf("parsing %s: %w", name, parseErr)
		}
		r.pages[name] = t
	}

	// Parse partial templates (standalone, no layout)
	partialFiles, err := fs.Glob(templateFS, "templates/partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("globbing partials: %w", err)
	}
	for _, path := range partialFiles {
		name := extractName(path)
		content, readErr := fs.ReadFile(templateFS, path)
		if readErr != nil {
			return nil, fmt.Errorf("reading %s: %w", path, readErr)
		}
		t, parseErr := template.New(name).Parse(string(content))
		if parseErr != nil {
			return nil, fmt.Errorf("parsing %s: %w", name, parseErr)
		}
		r.partials[name] = t
	}

	return r, nil
}

// RenderPage renders a full page template with the layout.
func (r *Renderer) RenderPage(w io.Writer, name string, data any) error {
	t, ok := r.pages[name]
	if !ok {
		return fmt.Errorf("page template %q not found", name)
	}
	return t.ExecuteTemplate(w, "layout", data)
}

// RenderPartial renders a partial template without layout.
func (r *Renderer) RenderPartial(w io.Writer, name string, data any) error {
	t, ok := r.partials[name]
	if !ok {
		return fmt.Errorf("partial template %q not found", name)
	}
	return t.Execute(w, data)
}

func extractName(path string) string {
	// "templates/pages/home.html" -> "home"
	base := ""
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			base = path[i+1:]
			break
		}
	}
	if base == "" {
		base = path
	}
	// Strip .html extension
	if len(base) > 5 && base[len(base)-5:] == ".html" {
		return base[:len(base)-5]
	}
	return base
}
