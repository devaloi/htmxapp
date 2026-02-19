# Build htmxapp — Server-Rendered Web App with Go, htmx, and templ

You are building a **portfolio project** for a Senior AI Engineer's public GitHub. It must be impressive, clean, and production-grade. Read these docs before writing any code:

1. **`G11-go-htmx-app.md`** — Complete project spec: architecture, phases, htmx patterns, commit plan. This is your primary blueprint. Follow it phase by phase.
2. **`github-portfolio.md`** — Portfolio goals and Definition of Done (Level 1 + Level 2). Understand the quality bar.
3. **`github-portfolio-checklist.md`** — Pre-publish checklist. Every item must pass before you're done.

---

## Instructions

### Read first, build second
Read all three docs completely before writing a single line of code. Understand the htmx interaction model (server returns HTML fragments, not JSON), the templ type-safe template system, the SSE real-time pattern, and the quality expectations.

### Follow the phases in order
The project spec has 6 phases. Do them in order:
1. **Project Foundation + Database** — Go project, SQLite, migrations, repository layer
2. **Auth + Middleware** — Session-based auth, CSRF, logging, recovery
3. **Templates + Layout** — templ base layout, auth pages, dashboard
4. **Kanban Board + htmx Interactions** — board page, task CRUD via htmx, drag-and-drop
5. **Real-Time Updates (SSE)** — event broker, live task updates across clients
6. **Polish + Testing + Documentation** — form validation, integration tests, README

### Commit frequently
Follow the commit plan in the spec. Use **conventional commits** (`feat:`, `test:`, `refactor:`, `docs:`, `chore:`). Each commit should be a logical unit.

### Quality non-negotiables
- **templ for all templates.** No `html/template` or `text/template`. All HTML generated via templ components with compile-time type checking.
- **htmx for all dynamic interactions.** No fetch/axios/XMLHttpRequest. The server returns HTML fragments. htmx swaps them into the DOM.
- **Session-based auth, not JWT.** This is a web app with cookies. JWT is for APIs. Use secure, httpOnly, sameSite cookies.
- **CSRF protection on all mutations.** Every POST/PUT/DELETE must include a valid CSRF token. Double-submit cookie pattern.
- **OOB swaps for toast notifications.** Every mutation response includes an out-of-band swap that adds a toast notification to the toast container.
- **Pure Go SQLite.** Use `modernc.org/sqlite` — no CGO, no external dependencies. Cross-platform builds.
- **Real SSE, not polling.** The event broker must use proper Server-Sent Events with goroutine-per-client, channel-based pub/sub.
- **Lint clean.** `golangci-lint run` must pass. `go vet` must pass.
- **Tests with goquery.** Parse HTML responses with goquery to assert on DOM structure, not string matching.
- **Minimal JavaScript.** The only JS should be htmx configuration and drag-and-drop (~50 lines max). Everything else is server-rendered.

### What NOT to do
- Don't use any JavaScript framework (React, Vue, Svelte, Alpine). The whole point is htmx.
- Don't return JSON from htmx endpoints. Return HTML fragments. This is the fundamental htmx pattern.
- Don't use `html/template`. Use templ exclusively for type safety.
- Don't use CGO-dependent SQLite drivers (mattn/go-sqlite3). Use `modernc.org/sqlite` for pure Go.
- Don't skip CSRF protection. Web apps without CSRF are vulnerable.
- Don't over-engineer the drag-and-drop. Keep it simple — the focus is on htmx server interactions, not fancy client-side UX.

---

## GitHub Username

The GitHub username is **devaloi**. For Go module paths, use `github.com/devaloi/htmxapp`. All internal imports must use this module path.

## Start

Read the three docs. Then begin Phase 1 from `G11-go-htmx-app.md`.
