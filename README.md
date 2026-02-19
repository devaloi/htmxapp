# htmxapp

A server-rendered contact manager built with Go and [htmx](https://htmx.org). No JavaScript frameworks — just HTML templates, htmx attributes, and a Go HTTP server.

## Features

- **Server-rendered HTML** — Go `html/template` with embedded templates
- **htmx interactions** — search, delete, and partial page updates without full reloads
- **Active search** — debounced search-as-you-type with `hx-trigger="input changed delay:300ms"`
- **CRUD operations** — create, read, update, delete contacts
- **Inline delete** — htmx DELETE swaps the row out of the DOM
- **Form validation** — server-side validation with error messages
- **Duplicate email detection** — prevents duplicate email addresses
- **Request logging** — structured logging with `slog`
- **Graceful shutdown** — clean shutdown on SIGINT/SIGTERM
- **Thread-safe store** — concurrent-safe in-memory storage with `sync.RWMutex`

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.26 |
| Server | `net/http` (stdlib) |
| Templates | `html/template` + `embed` |
| Interactivity | htmx 2.0 |
| Storage | In-memory (thread-safe) |
| Logging | `log/slog` |
| Styling | Custom CSS (no frameworks) |

## Architecture

```
htmxapp/
├── cmd/htmxapp/main.go            # Entry point
├── internal/
│   ├── handler/                    # HTTP handlers + middleware + static assets
│   │   ├── handler.go              # Routes and handler struct
│   │   ├── home.go                 # Home page
│   │   ├── contact.go              # Contact CRUD handlers
│   │   ├── middleware.go           # Logging, recovery, request ID
│   │   └── static/css/style.css    # Embedded stylesheet
│   ├── model/                      # Domain types
│   │   ├── contact.go              # Contact struct with validation
│   │   └── errors.go               # Domain errors
│   ├── server/                     # Server lifecycle
│   │   ├── server.go               # HTTP server with graceful shutdown
│   │   └── config.go               # Environment-based configuration
│   ├── store/                      # Data persistence
│   │   ├── store.go                # ContactStore interface
│   │   └── memory.go               # Thread-safe in-memory implementation
│   └── tmpl/                       # Template rendering
│       ├── render.go               # Template loader with embed.FS
│       └── templates/              # HTML templates
│           ├── layout.html         # Base layout
│           ├── pages/              # Full page templates
│           └── partials/           # htmx partial templates
```

## Prerequisites

- Go 1.22+ (uses new `net/http` routing patterns)

## Quick Start

```bash
# Clone and build
git clone https://github.com/devaloi/htmxapp.git
cd htmxapp
make build

# Run (seeds sample contacts)
make run
# → http://localhost:8080
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTMXAPP_HOST` | `""` | Bind address |
| `HTMXAPP_PORT` | `8080` | Listen port |
| `HTMXAPP_SEED` | `true` | Seed sample contacts on startup |

```bash
HTMXAPP_PORT=3000 HTMXAPP_SEED=false make run
```

## Development

```bash
make build    # Build binary
make run      # Build and run
make test     # Run tests with race detector
make lint     # Run golangci-lint
make fmt      # Format code
make clean    # Remove binary
```

## How htmx Works Here

### Search

The search input sends GET requests as the user types, replacing only the table body:

```html
<input type="search"
    hx-get="/contacts/search"
    hx-trigger="input changed delay:300ms"
    hx-target="#contact-rows">
```

### Delete

The delete button sends a DELETE request and removes the row from the DOM:

```html
<button
    hx-delete="/contacts/1"
    hx-target="#contact-1"
    hx-swap="outerHTML swap:200ms"
    hx-confirm="Delete Alice Johnson?">
```

### Server Response

For htmx requests, the server returns HTML partials instead of full pages. The `HX-Request` header distinguishes htmx requests from standard navigation.

## Tests

```bash
go test -race -v ./...
```

Tests cover:
- Model validation (valid/invalid contacts, email format)
- Store CRUD operations (create, read, update, delete, search)
- Store constraints (duplicate emails, not found errors)
- Concurrent access (parallel reads and writes)
- All HTTP handlers (happy path + error cases)
- Middleware (logging, recovery, request ID)
- Template rendering (pages, partials, missing templates)
- Server configuration (defaults, environment overrides)

## License

[MIT](LICENSE)
