# Cheap Eats

Find an affordable meal nearby, with prices for individual dishes. Cheap Eats is a Japan-focused app built around the question: **What can I eat nearby for the money I actually have?**

The current development version includes an Expo mobile/web client, a Go REST API, and PostgreSQL with PostGIS for geographic searches.

## Features

- Browse nearby dishes with prices in yen and Japanese/English display names.
- Filter by budget and search radius, sort by distance or price, and page through results.
- Search around Tokyo Station or use your device's location.
- View dish details, restaurant information, and price history.
- Create restaurants and dishes, and update prices through the API.

Maps, community price confirmations, and AI menu scanning/translation are planned features. See [project context](docs/project-context.md) for the broader product direction.

## Stack

| Component | Technology |
| --- | --- |
| Client | Expo, React Native, React, TypeScript, Expo Router |
| API | Go, Chi, `net/http` |
| Database | PostgreSQL 17, PostGIS 3.5, `pgx`, plain SQL |
| Local database | Docker Compose with a persistent volume |

## Run locally

### Prerequisites

- Go matching the version in [go.mod](go.mod) (currently `1.27.1`).
- Node.js and npm compatible with the Expo version in [mobile/package.json](mobile/package.json).
- Docker with Docker Compose.
- For native development: an iOS simulator or Android emulator, or a compatible device setup.

Run the following commands from the repository root unless stated otherwise.

### 1. Configure the database

Copy the example environment file if you do not already have a local `.env`:

```bash
cp .env.example .env
```

Set `POSTGRES_PASSWORD` and the matching password in `DATABASE_URL`. The defaults use database/user `cheap_eats` and host port **5433**. Keep local credentials out of version control.

```bash
make db-up
```

### 2. Apply migrations and load demo data

For a new database, load the environment and apply all upward migrations in order:

```bash
set -a
. ./.env
set +a

for migration in migrations/*.up.sql; do
  docker compose exec -T database \
    psql -X -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
    < "$migration" || break
done
```

Migrations are not applied automatically and this loop does not track applied versions. On an existing database, apply only migrations that have not already run. Resolve any migration error before continuing.

Load fictional bilingual dishes and restaurants around Tokyo Station:

```bash
docker compose exec -T database \
  psql -X -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
  < scripts/seed_demo.sql
```

The demo seed requires the database name `cheap_eats`. Its locations and prices are sample data, not verified restaurant listings. The smaller `scripts/seed.sql` fixture is used by integration tests.

### 3. Start the API

```bash
make run
```

The API listens at `http://localhost:8080`. `make run` starts the database if needed and loads `.env`; restart it after changing Go code.

In another terminal:

```bash
curl http://localhost:8080/api/health
curl 'http://localhost:8080/api/dishes/nearby?lat=35.6812&lng=139.7671&radius=2000&max_price=1000'
```

The health endpoint returns `{"status":"ok"}`; it does not check database readiness.

### 4. Start the client

Create or update `mobile/.env.local` with:

```dotenv
EXPO_PUBLIC_API_URL=http://localhost:8080
```

Then run:

```bash
cd mobile
npm ci
npm run web
```

For native development, use `npm run ios`, `npm run android`, or `npm start` to open the Expo launcher.

The API URL is the server origin, without `/api`. For the Android emulator, use `http://10.0.2.2:8080`; for a physical device, use your computer's reachable LAN address. Restart Expo after changing the environment file.

Web CORS currently allows read requests from `http://localhost:8081` and `http://127.0.0.1:8081`. If Expo uses a different web origin, update the CORS configuration in `cmd/api/main.go`.

The client initially searches around Tokyo Station. Use that preview area with the demo data, or choose **Use my location** to search nearby records.

## Development checks

Run backend unit and handler tests without a database:

```bash
make test
```

Check client types:

```bash
cd mobile
npx tsc --noEmit
```

For PostgreSQL integration tests, create a separate test database. From the repository root, with `.env` loaded as above:

```bash
docker compose exec -T database \
  createdb -U "$POSTGRES_USER" cheap_eats_test

for migration in migrations/*.up.sql; do
  docker compose exec -T database \
    psql -X -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d cheap_eats_test \
    < "$migration" || break
done

docker compose exec -T database \
  psql -X -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d cheap_eats_test \
  < scripts/seed.sql
```

Create and migrate the test database only once. Add `TEST_DATABASE_URL` to your root `.env`, using the same local credentials and port as `DATABASE_URL` but database name `cheap_eats_test`, then run:

```bash
make test-integration
```

GitHub Actions also runs backend unit/handler and PostgreSQL integration tests.

Other useful commands:

| Command | Purpose |
| --- | --- |
| `make fmt` | Format Go code |
| `make db-up` | Start the database and wait for its health check |
| `make db-down` | Stop the database while retaining its data volume |

## Repository layout

```text
cmd/api/             API entry point and routes
internal/handlers/   HTTP handlers and request validation
internal/repository/ Database queries and repository tests
internal/service/    Service logic
internal/models/     Shared backend data types
internal/database/   PostgreSQL connection setup
migrations/          SQL schema migrations
scripts/             Test fixtures and demo seed data
mobile/src/app/      Expo Router screens
mobile/src/api/      Client API calls and types
mobile/src/components/ UI components
docs/                API reference and product context
```

## API and project notes

See the [API guide](docs/api.md) for endpoints, request examples, validation rules, and response formats. Prices are integer yen amounts; distances are in meters and represent straight-line geographic distance.

The API currently has no authentication or authorization, including for write endpoints. It is a local-development MVP; public deployment needs access controls.

- [Product direction and architecture decisions](docs/project-context.md)
- [Architecture visualization](cheap-eats-architecture.html)
