# G11: htmxapp — Server-Rendered Web App with Go, htmx, and templ

**Catalog ID:** G11 | **Size:** M | **Language:** Go 1.26 + htmx 2.x + templ
**Repo name:** `htmxapp`
**One-liner:** A modern server-rendered web application using Go, htmx, and templ — a kanban-style task tracker with real-time updates, inline editing, and zero JavaScript frameworks.

---

## Why This Stands Out

- **Hypermedia-driven architecture** — demonstrates the htmx approach that's reshaping how developers think about web apps
- **templ for type-safe templates** — compile-time checked HTML generation, not text/template string soup
- **Server-Sent Events** — real-time task updates without WebSocket complexity
- **Zero JavaScript frameworks** — no React, Vue, or Svelte. Pure server-rendered HTML enhanced with htmx
- **Full CRUD with inline editing** — htmx swap targets, OOB swaps for toast notifications, optimistic UI patterns
- **Go standard library HTTP** — net/http with Go 1.22+ routing, middleware chain, clean architecture
- **Session-based auth** — secure cookies, CSRF protection, server-side sessions (not JWT for a web app)
- **Multi-model depth** — Users, Projects, Tasks with statuses, assignments, and real-time collaboration feel

---

## Architecture

```
htmxapp/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point: wire deps, run migrations, start server
├── internal/
│   ├── config/
│   │   └── config.go               # Env-based config with defaults
│   ├── database/
│   │   ├── database.go             # SQLite connection setup
│   │   ├── migrations.go           # Embedded SQL migration runner
│   │   └── migrations/
│   │       ├── 001_users.sql       # Users table
│   │       ├── 002_projects.sql    # Projects table
│   │       ├── 003_tasks.sql       # Tasks table with status, priority, assignee
│   │       └── 004_sessions.sql    # Sessions table for auth
│   ├── model/
│   │   ├── user.go                 # User struct + repository interface
│   │   ├── project.go              # Project struct + repository interface
│   │   ├── task.go                 # Task struct + status enum + repository interface
│   │   └── session.go              # Session struct + repository interface
│   ├── repository/
│   │   ├── user_repo.go            # SQLite user repository
│   │   ├── user_repo_test.go
│   │   ├── project_repo.go         # SQLite project repository
│   │   ├── project_repo_test.go
│   │   ├── task_repo.go            # SQLite task repository
│   │   ├── task_repo_test.go
│   │   └── session_repo.go         # SQLite session repository
│   ├── service/
│   │   ├── auth.go                 # Auth service: register, login, logout, session management
│   │   ├── auth_test.go
│   │   ├── project.go              # Project service: CRUD, membership
│   │   ├── project_test.go
│   │   ├── task.go                 # Task service: CRUD, status transitions, assignment
│   │   ├── task_test.go
│   │   └── events.go               # SSE event broker: subscribe, publish task updates
│   ├── handler/
│   │   ├── handler.go              # Base handler struct, route registration
│   │   ├── auth.go                 # Login page, register page, login/logout POST
│   │   ├── auth_test.go
│   │   ├── dashboard.go            # Dashboard page: list projects
│   │   ├── project.go              # Project pages: board view, settings
│   │   ├── project_test.go
│   │   ├── task.go                 # Task handlers: create, edit, move, delete (htmx partials)
│   │   ├── task_test.go
│   │   ├── events.go               # SSE endpoint: /events/tasks
│   │   └── response.go             # HTML response helpers, error rendering
│   ├── middleware/
│   │   ├── chain.go                # Middleware chaining helper
│   │   ├── logging.go              # Structured request logging
│   │   ├── auth.go                 # Session auth middleware (redirect to login)
│   │   ├── csrf.go                 # CSRF token generation and validation
│   │   ├── recovery.go             # Panic recovery
│   │   └── middleware_test.go
│   └── templates/
│       ├── layout.templ             # Base layout: head, nav, footer, toast container
│       ├── components/
│       │   ├── nav.templ            # Navigation bar with user menu
│       │   ├── toast.templ          # Toast notification (OOB swap target)
│       │   ├── modal.templ          # Modal dialog component
│       │   ├── form_error.templ     # Form field error display
│       │   └── pagination.templ     # Pagination controls
│       ├── pages/
│       │   ├── login.templ          # Login page
│       │   ├── register.templ       # Registration page
│       │   ├── dashboard.templ      # Dashboard: project cards
│       │   └── board.templ          # Kanban board: columns with task cards
│       └── partials/
│           ├── task_card.templ      # Single task card (htmx-enhanced)
│           ├── task_form.templ      # Task create/edit form (inline)
│           ├── task_column.templ    # Kanban column (droppable)
│           ├── project_card.templ   # Project summary card
│           └── project_form.templ   # Project create/edit form
├── static/
│   ├── css/
│   │   └── app.css                  # Custom styles (minimal, Tailwind handles most)
│   ├── js/
│   │   └── app.js                   # Minimal JS: htmx config, SSE reconnect, drag-and-drop
│   └── favicon.ico
├── go.mod
├── go.sum
├── Makefile
├── .env.example
├── .gitignore
├── .golangci.yml
├── LICENSE
└── README.md
```

---

## Page & Endpoint Reference

### Page Routes (return full HTML)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/` | No | Redirect to /dashboard or /login |
| `GET` | `/login` | No | Login page |
| `GET` | `/register` | No | Registration page |
| `GET` | `/dashboard` | Yes | Dashboard with project list |
| `GET` | `/projects/{id}` | Yes | Kanban board for project |
| `GET` | `/projects/new` | Yes | New project form |

### Action Routes (form submissions, return redirects or partials)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/login` | No | Process login form |
| `POST` | `/register` | No | Process registration form |
| `POST` | `/logout` | Yes | Logout, clear session |
| `POST` | `/projects` | Yes | Create project |
| `PUT` | `/projects/{id}` | Yes | Update project |
| `DELETE` | `/projects/{id}` | Yes | Delete project |

### htmx Partial Routes (return HTML fragments)

| Method | Path | Auth | htmx Target | Description |
|--------|------|------|-------------|-------------|
| `GET` | `/tasks/{id}/edit` | Yes | `#task-{id}` | Inline edit form for task |
| `POST` | `/projects/{pid}/tasks` | Yes | `#column-{status}` | Create task, append to column |
| `PUT` | `/tasks/{id}` | Yes | `#task-{id}` | Update task, swap card |
| `PATCH` | `/tasks/{id}/status` | Yes | `#task-{id}` | Move task to new status column |
| `DELETE` | `/tasks/{id}` | Yes | `#task-{id}` | Remove task card |

### SSE Events

| Event | Path | Data | Description |
|-------|------|------|-------------|
| `task-created` | `/events/tasks?project={id}` | HTML fragment | New task card to append |
| `task-updated` | `/events/tasks?project={id}` | HTML fragment | Updated task card to swap |
| `task-deleted` | `/events/tasks?project={id}` | Task ID | Task ID to remove |
| `task-moved` | `/events/tasks?project={id}` | HTML fragment | Task card for column swap |

### htmx Patterns Used

| Pattern | Where | Description |
|---------|-------|-------------|
| `hx-get` + `hx-target` | Task card click | Load inline edit form into card slot |
| `hx-post` + `hx-swap="beforeend"` | New task form | Append new card to column |
| `hx-put` + `hx-swap="outerHTML"` | Edit task form | Replace card with updated version |
| `hx-delete` + `hx-swap="delete"` | Delete button | Remove card from DOM |
| OOB swap | Toast notifications | `hx-swap-oob="afterbegin:#toast-container"` |
| `hx-trigger="sse:task-updated"` | Kanban board | Listen to SSE for real-time updates |
| `hx-confirm` | Delete actions | Browser confirm dialog before destructive actions |
| `hx-indicator` | All forms | Show loading spinner during requests |

---

## Task Status Flow

```
┌──────────┐    ┌─────────────┐    ┌────────┐    ┌──────┐
│ BACKLOG  │ →  │ IN_PROGRESS │ →  │ REVIEW │ →  │ DONE │
└──────────┘    └─────────────┘    └────────┘    └──────┘
     ↑               ↑                 ↑
     └───────────────┴─────────────────┘
              (can move backwards)
```

---

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.26 |
| HTTP | net/http (Go 1.22+ routing) |
| Templates | templ (type-safe HTML generation) |
| Interactivity | htmx 2.x (hypermedia-driven) |
| Real-time | Server-Sent Events (SSE) |
| Database | SQLite 3 via `modernc.org/sqlite` (pure Go) |
| Migrations | Embedded SQL files (embed package) |
| Styling | Tailwind CSS (via CDN) |
| Auth | Session-based (secure cookies + server-side sessions) |
| CSRF | Double-submit cookie pattern |
| Testing | stdlib testing + httptest + goquery |
| Linting | golangci-lint |
| Drag & Drop | Minimal vanilla JS (SortableJS or ~30 lines custom) |

---

## Phased Build Plan

### Phase 1: Project Foundation + Database

**1.1 — Project setup**
- `go mod init github.com/devaloi/htmxapp`
- Install templ: `go install github.com/a-h/templ/cmd/templ@latest`
- Add dependencies: modernc.org/sqlite, github.com/a-h/templ
- Create directory structure, Makefile, .gitignore, .golangci.yml
- Makefile targets: build, test, lint, run, generate (templ), dev (with air live-reload)

**1.2 — Database and migrations**
- SQLite connection with WAL mode and busy timeout
- Embedded migration runner using `embed` package
- 001: users (id, email, password_hash, name, created_at)
- 002: projects (id, name, description, owner_id FK, created_at, updated_at)
- 003: tasks (id, title, description, status, priority, project_id FK, assignee_id FK, position, created_at, updated_at)
- 004: sessions (token, user_id FK, expires_at, created_at)
- Tests: migrations run, tables created

**1.3 — Model and repository layer**
- User, Project, Task, Session structs
- Task status enum: BACKLOG, IN_PROGRESS, REVIEW, DONE
- Task priority enum: LOW, MEDIUM, HIGH
- Repository interfaces + SQLite implementations
- CRUD operations for each model
- Tests: full CRUD for each repository

### Phase 2: Auth + Middleware

**2.1 — Auth service**
- Register: validate input, hash password (bcrypt), create user
- Login: verify credentials, create session, set secure cookie
- Logout: delete session, clear cookie
- Session lookup: get user from session token
- Tests: register, login, logout, expired session

**2.2 — Middleware stack**
- Logging middleware: method, path, status, duration
- Auth middleware: check session cookie, load user into context, redirect to /login if missing
- CSRF middleware: generate token on GET, validate on POST/PUT/DELETE
- Recovery middleware: catch panics, render 500 page
- Tests: auth redirects, CSRF validation, panic recovery

### Phase 3: Templates + Layout

**3.1 — templ setup and base layout**
- Install templ, configure generation in Makefile
- Base layout: HTML head (Tailwind CDN, htmx CDN), nav, main content area, toast container, footer
- Navigation component: logo, project links, user menu with logout
- Toast component: success/error/info variants, auto-dismiss, OOB swap target
- Static file serving: /static/ path

**3.2 — Auth pages**
- Login page: email + password form, error display, link to register
- Register page: name + email + password form, validation errors
- Forms submit via standard POST (no htmx for auth — full page flow)
- CSRF token hidden field in all forms
- Tests: pages render, forms contain CSRF tokens

**3.3 — Dashboard page**
- List user's projects as cards
- Each card: name, description, task count by status, link to board
- "New Project" button → modal or separate page
- Project create/edit forms
- Tests: dashboard renders with projects

### Phase 4: Kanban Board + htmx Interactions

**4.1 — Board page**
- Four columns: Backlog, In Progress, Review, Done
- Each column renders task cards sorted by position
- Task card: title, priority badge, assignee avatar, edit/delete buttons
- Full page load fetches all tasks for project

**4.2 — Task CRUD via htmx**
- Create task: form at top of column, `hx-post`, new card appended to column
- Edit task: click card → `hx-get` loads inline edit form → `hx-put` swaps back to card
- Delete task: `hx-delete` with `hx-confirm`, card removed from DOM
- Move task: `hx-patch` changes status, card moves to new column
- Toast notification on every action via OOB swap
- Tests: each htmx endpoint returns correct HTML fragment

**4.3 — Drag and drop (minimal JS)**
- Minimal vanilla JS or SortableJS for drag-and-drop between columns
- On drop: `hx-patch` to update task status and position
- Reorder within column: update position field
- Keep JS under 50 lines if hand-rolled

### Phase 5: Real-Time Updates (SSE)

**5.1 — SSE event broker**
- In-memory broker: subscribe (channel per project), publish, unsubscribe
- Goroutine-safe with mutex or channels
- Client connects to `/events/tasks?project={id}`
- Server pushes events: task-created, task-updated, task-deleted, task-moved
- Events contain HTML fragments (htmx can swap directly)
- Auto-reconnect on client side (htmx SSE extension)
- Tests: broker subscribe/publish/unsubscribe, concurrent access

**5.2 — Integrate SSE with task operations**
- On task create → publish task-created event with rendered card HTML
- On task update → publish task-updated event with re-rendered card
- On task delete → publish task-deleted event with task ID
- On task move → publish task-moved event with card HTML + old/new column IDs
- Board page listens to SSE and swaps fragments in real-time
- Tests: create task triggers SSE event

### Phase 6: Polish + Testing + Documentation

**6.1 — Form validation**
- Server-side validation for all forms
- Render validation errors inline next to form fields
- Use templ components for consistent error display
- Required fields, length limits, email format, unique constraints
- Tests: validation errors render correctly

**6.2 — Integration tests**
- Use httptest for full request/response testing
- Use goquery to parse HTML responses and assert content
- Test full flows: register → login → create project → create tasks → move tasks → logout
- Test auth protection: unauthenticated access redirected
- Test CSRF: requests without valid token rejected
- Test htmx partials: correct HTML fragments returned

**6.3 — Error pages**
- Custom 404 page (styled, link back to dashboard)
- Custom 500 page (styled, generic error message)
- Consistent with site layout

**6.4 — README**
- Badges (CI, Go version, license)
- Screenshot or GIF of the kanban board in action
- Quick start: go run, open browser
- Architecture overview: why htmx over SPA, why templ over text/template
- Feature list with htmx patterns used
- Development commands (Makefile targets)
- Tech stack rationale

---

## Commit Plan

1. `chore: scaffold Go project with dependencies and Makefile`
2. `feat: add SQLite database with embedded migrations`
3. `feat: add model structs and repository layer`
4. `feat: add auth service with session-based login`
5. `feat: add middleware stack (logging, auth, CSRF, recovery)`
6. `feat: add templ base layout and navigation`
7. `feat: add auth pages (login, register)`
8. `feat: add dashboard page with project list`
9. `feat: add project CRUD handlers`
10. `feat: add kanban board page with task columns`
11. `feat: add task CRUD with htmx partials`
12. `feat: add drag-and-drop task movement`
13. `feat: add toast notifications via OOB swaps`
14. `feat: add SSE event broker for real-time updates`
15. `feat: integrate SSE with task operations`
16. `feat: add form validation with inline errors`
17. `test: add integration tests with httptest and goquery`
18. `feat: add custom error pages (404, 500)`
19. `docs: add README with screenshots and architecture overview`
