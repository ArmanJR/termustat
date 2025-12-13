# Termustat

This is the repository for **Termustat**, which is an online course timetabling / planning app for university students.

## What it does

- Lets students build a weekly schedule on a visual calendar, rearrange classes, and spot time conflicts
- Provides search/filtering by things like faculty, professor, semester, and keywords
- Imports/exports course data via a dedicated **engine** that parses "Golestan" exports (e.g., raw HTML/SQL) into structured data

## How it’s built (high level)

- Engine (Go): parser/transformer for Golestan export files
- API (Go): REST backend for courses/faculties/semesters/users, plus auth + email workflows
- Frontend (React): interactive calendar UI + course selection/filtering
- PostgreSQL as the database, with Docker/Docker Compose to run the whole stack

## Where the code lives

The remote repository is publicly available on GitHub: `https://github.com/ArmanJR/termustat`.

You can use GitHub CLI `gh` when needed.

## Development

### What “good” looks like

* Ship in **small, reviewable increments** (prefer multiple PRs over one giant PR).
* Optimize for **clarity over cleverness**: readable code, explicit names, predictable structure.
* Keep **contracts stable** (API + Engine payloads + DB schema) and change them deliberately with versioning and migration plans.
* “Fix forward”: when you touch an area, **leave it better** (tests, types, logs).

### Ownership & boundaries

* The repo is one module with three components:

  * `api/` (Go/Gin backend)
  * `engine/` (Go/Gin parser service)
  * `frontend/` (React SPA)
* Each PR should clearly state which component(s) it affects and what the integration impact is (routes, DTOs, env vars, compose).

---

## Local development workflow

### Standard way to run

* Use Docker Compose for local dev to match prod routing assumptions (Nginx reverse proxy → API/Engine + Frontend).
* Prefer running everything through the gateway so paths behave like prod:

  * UI via gateway
  * API under `/api/...`
  * Engine under `/engine/...`

### Environment & secrets

* Never commit secrets. Use:

  * `.env` locally
  * `.env.example` in repo (safe placeholders)
* If you add a required env var:

  * Update `.env.example`
  * Update docs (README / this doc)
  * Fail fast on boot with a clear error message

---

## Branching, commits, and PR hygiene

### Branching

* `main` is always releasable.
* Use short-lived feature branches:

  * `feat/<topic>`
  * `fix/<bug>`
  * `refactor/<area>`

### Commits

* Prefer meaningful, scoped commits.
* Suggested convention:

  * `api: ...`, `engine: ...`, `frontend: ...`, `infra: ...`, `docs: ...`

### PR requirements (checklist)

* Clear description: **problem → approach → tradeoffs**
* Includes:

  * tests (unit/integration) or explanation why not
  * docs updates if behavior changed (routes/DTO/env vars)
* No debug logging, no commented-out code, no “TODO” without an issue link.
* Reviewer-friendly:

  * keep PRs under ~300–500 LOC when possible
  * avoid mixing unrelated refactors with feature changes

---

## Code organization & style

### General

* Keep dependency direction:

  * handlers/controllers → services → repositories → models
* Don’t let HTTP concerns leak into service logic:

  * service methods take typed inputs and return typed outputs/errors
  * handlers map HTTP → DTO → service and service result → HTTP response

### Go guidelines (API + Engine)

* Use `context.Context` (from Gin request) for:

  * DB calls
  * external calls (mail)
  * timeouts/cancellation
* Error handling:

  * wrap errors with context (`fmt.Errorf("...: %w", err)`)
  * never swallow errors
* Logging:

  * log at boundaries (HTTP entry/exit, external calls, DB failures)
  * include request correlation info if available (request id, user id)
* Prefer explicit DTOs at the boundary:

  * request/response DTOs are stable contracts
  * internal models can evolve

### Frontend guidelines

* Treat API responses as contracts:

  * validate/guard against missing fields
  * centralize API client logic (base URL, auth headers, refresh flow)
* Keep UI state predictable:

  * avoid “state soup” scattered across components
  * favor a small number of well-owned state domains (auth, timetable, catalog)

---

## API design rules (backend)

### Request/response conventions

* Define a consistent error shape (recommended).
* Always return stable identifiers and machine-readable error codes.
* Don’t return raw DB/internal errors to clients.

### Auth & security

* JWT claims and scopes:

  * keep authorization checks centralized (middleware + service guardrails)
* Passwords:

  * always hashed with a strong algorithm
  * never logged
* Emails:

  * verification/reset links must be derived from configured frontend base URL
* Rate limiting & brute-force protection should exist (at least on login/reset endpoints) before production.

---

## Engine rules (Golestan parsing)

### Contract stability

* The engine’s `/process` interface is an integration contract:

  * request form field names
  * response JSON schema
* If you need to change it:

  * add fields without breaking existing ones
  * update the API import adapter in the same PR

### Parser robustness

* Assume real-world HTML is messy:

  * tolerate missing cells/rows
  * be explicit about defaults and error messages
* Add fixture tests:

  * store anonymized HTML samples
  * assert parsed output (golden tests)

---

## Database & schema management

### Schema changes

* Every schema change must be:

  * explicit (migration)
  * reversible when feasible
  * documented (why it changed, what data is impacted)
* Never rely on “it works on my machine” auto-migrations for production.

### Data integrity rules

* Enforce invariants close to the data:

  * unique constraints (e.g., course code within semester)
  * foreign keys where appropriate
* Domain rules belong in services too (capacity/gender restriction/time overlap), but DB should still protect critical invariants.

---

## Testing strategy

### Minimum expectations

* Backend:

  * unit tests for service logic (conflict detection, validation rules)
  * repository tests where queries are non-trivial
  * one or two integration tests that boot the API and hit key routes
* Engine:

  * golden tests for parsing
* Frontend:

  * component tests for timetable rendering + conflict UI states
  * lightweight e2e smoke test if possible (login → select course → validate)

### CI gate (recommended)

* Format + lint + test must pass before merge.
* Block merges if:

  * tests fail
  * swagger/docs are out of sync with routes/contracts

---

## Observability & operability

### Logging

* Structured logs (JSON if possible in containers).
* Log key events:

  * auth success/failure (without sensitive info)
  * course select/unselect
  * import runs and counts (how many records accepted/rejected)
* Include identifiers:

  * request id, user id, semester id, course id when available

### Health & readiness

* Keep `/health` fast and dependency-aware:

  * DB reachability for readiness
  * process liveness for health

---

## Definition of Done

A change is “done” when:

* It’s deployed locally via compose and exercised through the gateway paths (`/api`, `/engine`)
* It has tests or a documented reason it doesn’t
* It updates docs/Swagger when contracts changed
* It doesn’t introduce breaking changes without a versioning plan
* It includes a rollback/mitigation note if risk is non-trivial

---

## Notes

* Prefer “proposal PRs” early for big refactors (skeleton + plan + small first step).
* When you discover a flaw (inconsistency, missing validation, shaky route naming), either fix it or file an issue with a concrete reproduction and suggested direction.
