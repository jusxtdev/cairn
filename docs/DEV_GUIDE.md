# Development Guide

How to run Cairn locally: Postgres in Docker, the API with live reload.

## Prerequisites

- **Go 1.25+** — `go version`
- **Docker** with the Compose plugin — `docker compose version`
- **Air** (optional, for live reload) — see step 3
- **golang-migrate** CLI (for database migrations) — see step 4

## 1. Configure environment

From the repository root (`cairn/`):

```bash
cp .env.example .env
```

Open `.env` and set your local values:

```env
POSTGRES_USER=cairn
POSTGRES_PASSWORD=change-me
POSTGRES_DB=cairn
```

`.env` is gitignored — never commit it. `.env.example` is the committed template.

> `docker-compose.yml` requires `POSTGRES_PASSWORD`. If `.env` is missing, Compose
> fails with `POSTGRES_PASSWORD is required - copy .env.example to .env`.

## 2. Start the local database

```bash
docker compose up -d db
```

Check that it is running and healthy:

```bash
docker compose ps
docker compose logs -f db
```

Connection details :

| Setting  | Value                          |
| -------- | ------------------------------ |
| Host     | `localhost`                    |
| Port     | `5433`                         |
| Database | value of `POSTGRES_DB`         |
| User     | value of `POSTGRES_USER`       |
| Password | value of `POSTGRES_PASSWORD`   |

Connect through the container (no password prompt needed):

```bash
docker compose exec db psql -U cairn -d cairn
```

Or from a host-installed `psql`, using the values from your `.env`:

```bash
psql "postgresql://<user>:<password>@localhost:5433/<db>"
```

> `POSTGRES_*` variables are only applied the first time the database volume is
> initialized. Changing the password later does not update an existing volume.
> Either run `docker compose down -v` to reset the data, or change it in SQL with
> `ALTER USER cairn WITH PASSWORD 'new-password';`.

### Useful Docker commands

```bash
docker compose stop db      # stop, keep data
docker compose down         # remove containers, keep data volume
docker compose down -v      # remove containers AND the database volume (data loss)
```

## 3. Run the API

### With live reload (Air)

Install Air once:

```bash
go install github.com/air-verse/air@latest
```

Make sure the Go bin directory is on your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then run **from the repository root** (`cairn/`) — not from `cmd/api`, because
Air reads `.air.toml` from the current directory:

```bash
air
```

Air watches `.go` files, rebuilds into `./tmp/main`, and restarts the server on
every change. Build failures are written to `tmp/build-errors.log`.

### Without live reload

```bash
go run ./cmd/api
```

### Verify it works

The API listens on `http://localhost:8080`:

```bash
curl -s http://localhost:8080/health
```

Expected response:

```json
{"status":true,"message":"API is healthy"}
```

Stop the server with `Ctrl+C`.

## 4. Database migrations

Migrations are managed with [golang-migrate](https://github.com/golang-migrate/migrate).
Files live in `migrations/` and follow the naming convention
`<version>_<name>.up.sql` / `<version>_<name>.down.sql`.

### Install the CLI

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

On macOS you can use `brew install golang-migrate` instead. Verify with
`migrate -version`.

### Set the database URL

The CLI does not read `.env` automatically, so load it into your shell first:

```bash
set -a
source .env
set +a

export DATABASE_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:5433/$POSTGRES_DB?sslmode=disable"
```

> If your password contains URL-special characters (`@`, `:`, `/`, `#`, ...),
> URL-encode them in `DATABASE_URL`.

### Apply migrations

Run from the repository root:

```bash
migrate -path migrations -database "$DATABASE_URL" up      # apply all pending
migrate -path migrations -database "$DATABASE_URL" down 1  # roll back one
migrate -path migrations -database "$DATABASE_URL" version # current version
```

`migrate` tracks the applied version and dirty state in a `schema_migrations`
table inside the database.

### Create a new migration

```bash
migrate create -ext sql -dir migrations -seq add_users_table
```

This generates the next numbered `.up.sql` and `.down.sql` pair. Always write
both directions: `up` applies the change, `down` reverses it.

### Important

- Never edit a migration that has already been applied — add a new one instead.
- If a migration fails halfway, the database is marked **dirty** and `migrate`
  refuses to continue. Fix the cause, then force the version you are actually
  on (the version before the failed one) before retrying:

  ```bash
  migrate -path migrations -database "$DATABASE_URL" force <version>
  migrate -path migrations -database "$DATABASE_URL" up
  ```

  Only use `force` after you have manually verified the schema state.

## Project layout

```
cairn/
├── cmd/api/            # application entrypoint (main.go)
├── internal/           # application packages (httpx, user, ...)
├── migrations/         # golang-migrate SQL files (up/down)
├── tmp/                # Air build output (gitignored)
├── .air.toml           # Air live-reload config
├── .env                # local secrets (gitignored)
├── .env.example        # committed template for .env
└── docker-compose.yml  # local Postgres service
```

## Troubleshooting

- **`air: command not found`** — add `$(go env GOPATH)/bin` to `PATH`, then
  reopen the shell.
- **`POSTGRES_PASSWORD is required`** — create `.env` from `.env.example`.
- **Port 5433 already in use** — stop the other Postgres, or change the host
  port in the `ports` mapping of `docker-compose.yml`.
- **Port 8080 already in use** — stop whatever is bound to it, or change the
  `:8080` address in `cmd/api/main.go`.
- **Stale build output** — stop Air and delete `tmp/` before restarting.
- **Password change has no effect** — the Postgres volume was already
  initialized; see the note in step 2.
- **`migrate: command not found`** — install the CLI and add
  `$(go env GOPATH)/bin` to `PATH`.
- **`Dirty database version N`** — a migration failed partway. Fix it, then
  `force` the last good version; see step 4.
- **`error: no change`** — the database is already up to date; nothing to do.
