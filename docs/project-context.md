# Cheap Eats project context

This document preserves the important decisions from earlier Cheap Eats conversations so future work in this repository can continue from the same context.

## Product direction

Cheap Eats is a dish-first discovery app for finding a specific affordable meal nearby. Its core question is:

> What can I eat nearby for the money I actually have?

The initial market is Japan, with a focused launch in a dense Tokyo neighborhood and an initial promise centered on meals under ¥1,000. The product can later support other budgets, cities, and currencies.

Restaurants are supporting entities. Search results should prioritize individual dishes with an exact price, rather than restaurant-level price symbols.

## Initial user experience

- Browse affordable nearby dishes in map and list views.
- Filter by maximum price, distance, cuisine, dietary needs, and whether a restaurant is open.
- View the dish, current price, restaurant, location, photographs, and the last price-verification date.
- Submit dishes, photographs, and corrected prices.
- Confirm that an existing price is still accurate or report it as outdated.
- Open a restaurant page to see its details and all known affordable dishes.

Accounts, social feeds, delivery integrations, and sophisticated recommendations are later-stage features unless the MVP proves that they are required.

## AI menu translation

A user can photograph a Japanese menu. The system should detect dishes and prices, extract and translate the Japanese text, and reconstruct a readable translated menu while preserving the relationship between each dish and its price. Extracted data should be reviewed by a person before it becomes trusted public data.

## Data model

The core relationship is:

```text
Restaurant
└── Dishes
    ├── current price and price history
    ├── category and photographs
    ├── price confirmations
    └── reviews or comments
```

Likely entities include users, restaurants, locations, dishes, prices, price history, menus, photographs, contributions, verifications, favorites, reviews, promotions, and business accounts.

## Technical direction

- Mobile app: React Native and TypeScript.
- Backend API: Go, using the standard `net/http` package with Chi.
- Architecture: a modular monolith initially, not microservices.
- Primary database: PostgreSQL with PostGIS for distance and map queries.
- Database access: `pgx` and plain SQL initially; consider `sqlc` later.
- Object storage: menu and dish photographs.
- Maps and place data: Google Maps SDK and Google Places API were discussed; provider cost and data-use terms should be checked before implementation.
- Redis: introduce only when caching, rate limiting, sessions, or background-job coordination create a clear need.
- AI processing: a multimodal vision model for OCR, structured menu extraction, translation, categorization, and moderation.
- Deployment: Docker and a managed container platform, keeping the system portable.

An early API shape discussed was:

```text
GET  /api/dishes?lat=...&lng=...&max_price=1000
GET  /api/dishes/{id}
POST /api/dishes
POST /api/dishes/{id}/confirm-price
```

Suggested backend structure:

```text
cmd/api/main.go
internal/handlers/
internal/service/
internal/repository/
internal/models/
migrations/
```

## Monetization direction

Consumer discovery should remain free. The strongest early revenue paths are:

1. Clearly labeled promoted dishes that still match a user's budget and location.
2. Restaurant Pro subscriptions for claimed profiles, menu and price management, translated menus, analytics, promotion scheduling, and verification.
3. Paid multilingual menu conversion, offered once or as part of Pro.
4. Coupons, reservations, takeaway, prepaid-meal, or delivery-partner referral fees.
5. Aggregated, privacy-safe restaurant analytics.

## Recommended delivery order

1. Create and understand the Go project structure.
2. Start a basic HTTP server and `/api/health` endpoint.
3. Define restaurant and dish models.
4. Build temporary in-memory dish endpoints.
5. Add PostgreSQL and PostGIS.
6. Build the React Native client.
7. Add maps and location-based search.
8. Seed one launch neighborhood with trustworthy dish and price data.
9. Add community price verification.
10. Add menu scanning and translation.

## Collaboration preference

Development should be guided and educational: provide the code, explain which files to create, describe each file's responsibility, and explain the meaning of important code. Do not modify application files unless explicitly requested.

## Imported conversation references

- `Start project with guided coding` — Codex task `01a0857e-bcee-7d03-82fb-be23ac4ba8d1`
- `Summarize Cheap Eats stack` — Codex task `01a08324-1b9a-75b3-ba13-ec0d7b1642ce`
- `Import Cheap Eats Conversations` — ChatGPT conversation `6aa08c2b-9c18-83e8-a13e-fcb844844dbc`

