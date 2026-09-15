# Cheap Eats API guide

This guide describes the current local-development API, not planned features.
Examples use fictional data. Generated IDs, prices, and timestamps will vary.

## Local setup

With your local `.env` and database migrations already configured, start the API
from the project directory:

```bash
make run
```

Base URL: `http://localhost:8080`. Restart the API after changing Go code.
Never commit `.env` or copy its credentials into this guide.

## Conventions and limitations

- Successful responses use `application/json`. Lists are JSON arrays; empty lists
  are `[]`. There is no response envelope or total-count field.
- Errors are currently plain text, not JSON, and normally end with a newline.
- IDs in paths must be positive integers that fit in a signed 64-bit integer.
- Send JSON bodies with `Content-Type: application/json`. The current handlers
  decode JSON but do not enforce this header or return `415` for a wrong header.
- JSON body endpoints accept at most 64 KiB and reject unknown fields, malformed
  JSON, and extra JSON values. Duplicate JSON field names are not rejected.
- Prices are integers. For the Japan-focused examples, `950` means 950 yen.
  Creating or updating a price accepts 1 through 2,147,483,647 inclusive.
- Coordinates use latitude/longitude in degrees. Latitude must be between -90
  and 90, longitude between -180 and 180; non-finite values are invalid.
- Distances and search radii are in meters. Distances are straight-line geographic
  distances, not walking routes or travel times.
- There is currently no authentication or authorization. Anyone who can reach
  the API can read data, create restaurants/dishes, and update prices. Do not
  expose these write endpoints publicly before adding access controls.
- Price filtering and sorting compare stored numbers without currency filtering
  or conversion. Use JPY data for meaningful comparisons in the current MVP.

## Endpoint overview

| Method | Path | Successful response |
| --- | --- | --- |
| GET | `/api/health` | `200`, health object |
| GET | `/api/dishes` | `200`, array of dishes |
| GET | `/api/dishes/nearby` | `200`, array of nearby dish results |
| GET | `/api/dishes/{id}` | `200`, dish object |
| POST | `/api/dishes` | `201`, created dish object |
| PATCH | `/api/dishes/{id}/price` | `200`, updated dish object |
| GET | `/api/dishes/{id}/price-history` | `200`, array of history entries |
| GET | `/api/restaurants` | `200`, array of restaurants |
| GET | `/api/restaurants/nearby` | `200`, array of restaurant suggestions |
| GET | `/api/restaurants/{id}` | `200`, restaurant object |
| POST | `/api/restaurants` | `201`, created restaurant object |
| GET | `/api/restaurants/{id}/dishes` | `200`, array of dishes |

## Health

```bash
curl -i http://localhost:8080/api/health
```

Response: `200 OK` with `{"status":"ok"}`.
This handler does not query PostgreSQL; it is not a database-readiness check.

## Dish response shape

```json
{
  "id": 1,
  "name": "Shoyu Ramen",
  "price": 850,
  "currency": "JPY",
  "restaurant_id": 1
}
```

### List dishes

```bash
curl -i http://localhost:8080/api/dishes
curl -i "http://localhost:8080/api/dishes?max_price=1000"
```

Optional `max_price` must be a positive integer. The filter is inclusive.
Omitting it, or supplying an empty value, returns all dishes.
PostgreSQL returns unfiltered dishes in ID order and filtered dishes in price,
then ID order. This endpoint is not paginated.

### Get one dish

```bash
curl -i http://localhost:8080/api/dishes/1
```

Returns `400` for an invalid ID or `404` with `dish not found` for a missing dish.

### Create a dish

This example writes to your development database. Replace `restaurant_id` with
an existing restaurant's ID. Each successful request creates a new dish.

```bash
curl -i http://localhost:8080/api/dishes \
  -H "Content-Type: application/json" \
  -d '{"restaurant_id":1,"name":"Chicken Curry","price":900,"currency":"JPY"}'
```

Required fields: nonblank `name`, positive `restaurant_id`, and valid `price`.
Names are trimmed. Currency is trimmed and uppercased; omitted or blank currency
defaults to `JPY`. A missing referenced restaurant returns `400`, not `404`.
Success returns `201` with the created dish, including its generated ID.

Dish creation accepts only JPY after trimming and uppercasing. Omitted or blank currency defaults to JPY. Other values return 400 with currency must be JPY.

### Update a dish's price

This changes the stored price of the selected dish:

```bash
curl -i -X PATCH http://localhost:8080/api/dishes/1/price \
  -H "Content-Type: application/json" \
  -d '{"price":950}'
```

Only `price` is accepted. Missing, null, non-integer, or out-of-range prices are
rejected with `400`. Missing dishes return `404`.

Success returns `200` with the updated dish. Its name, restaurant, and currency
remain unchanged. PostgreSQL saves the price update and history entry in one
transaction. Repeating the current price succeeds without another history entry
or a change to `updated_at`.

### Search nearby dishes

```bash
curl -i -G http://localhost:8080/api/dishes/nearby \
  --data-urlencode "lat=35.6812" \
  --data-urlencode "lng=139.7671" \
  --data-urlencode "radius=2000" \
  --data-urlencode "max_price=1000" \
  --data-urlencode "sort=price" \
  --data-urlencode "limit=20" \
  --data-urlencode "offset=0"
```

| Parameter | Required | Rules / default |
| --- | --- | --- |
| `lat` | Yes | Latitude, -90 through 90 |
| `lng` | Yes | Longitude, -180 through 180 |
| `max_price` | Yes | Positive integer; inclusive ceiling |
| `radius` | No | 1 through 5000 meters; default 1000 |
| `sort` | No | `distance` (default) or `price` |
| `limit` | No | Integer 1 through 100; default 20 |
| `offset` | No | Nonnegative integer; default 0 |

Empty `radius` and `sort` values use their defaults. Empty `limit` or `offset`
values are invalid. Invalid parameters return `400`.

Example response:

```json
[
  {
    "dish": {
      "id": 1,
      "name": "Shoyu Ramen",
      "price": 850,
      "currency": "JPY",
      "restaurant_id": 1
    },
    "restaurant": {
      "id": 1,
      "name": "Tokyo Ramen",
      "address": "Marunouchi, Tokyo",
      "latitude": 35.6812,
      "longitude": 139.7671
    },
    "distance_meters": 0
  }
]
```

`distance` sorts by actual distance, then price, then dish ID. `price` sorts by
price, then actual distance, then dish ID. Returned distances are rounded to
whole meters, so displayed distances can tie even when actual distances differ.
No matches return `200` with `[]`.

### Read price history

```bash
curl -i "http://localhost:8080/api/dishes/1/price-history?limit=20&offset=0"
```

`limit` defaults to 20 and accepts integers 1 through 100. `offset` defaults to 0
and accepts nonnegative integers. Explicitly empty values are invalid.

Example response:

```json
[
  {
    "id": 1,
    "dish_id": 1,
    "old_price": 850,
    "new_price": 950,
    "currency": "JPY",
    "changed_at": "2026-09-16T10:00:00Z"
  }
]
```

Entries are returned in descending history ID order. Existing dishes with no
recorded changes, or offsets past the final entry, return `200` with `[]`.
Invalid IDs/pagination return `400`; missing dishes return `404`.

Only changes made after history recording was implemented are available. Initial
dish creation does not add a history entry. Deleting a dish also deletes its
history. No submitter identity is currently stored or returned.

For the next page, increase `offset` by `limit`. There is no total count or next
page token. New changes between requests can shift offset pages and repeat
entries. Pagination limits response size; it does not restrict public access.

## Restaurant response shape

```json
{
  "id": 1,
  "name": "Tokyo Ramen",
  "address": "Marunouchi, Tokyo",
  "latitude": 35.6812,
  "longitude": 139.7671
}
```

### List or get restaurants

```bash
curl -i http://localhost:8080/api/restaurants
curl -i http://localhost:8080/api/restaurants/1
```

The list is unpaginated and ordered by ID. Getting one restaurant returns `400`
for an invalid ID or `404` with `restaurant not found` for a missing restaurant.

### Create a restaurant

This writes a fictional restaurant to your development database. Each successful
request creates another record; duplicate detection is not implemented.

```bash
curl -i http://localhost:8080/api/restaurants \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Curry Kitchen","address":"Test address, Tokyo","latitude":35.6830,"longitude":139.7690}'
```

All four fields are required. Name and address are trimmed and must not be blank.
Coordinates must be JSON numbers in the valid ranges; zero is allowed, while
missing or null coordinates are rejected. Success returns `201` with the created
restaurant and a `Location` header such as `/api/restaurants/3`.

### Search nearby restaurants

```bash
curl -i "http://localhost:8080/api/restaurants/nearby?lat=35.6812&lng=139.7671&radius=2000"
```

`lat` and `lng` are required, with the usual coordinate ranges. Optional `radius`
accepts 1 through 5000 meters and defaults to **500** when omitted or empty
(unlike the nearby-dishes default of 1000).

Returns an array of restaurant objects with an additional `distance_meters`
field at the same level as `id` and `name`, not a nested `restaurant` object.
Results are ordered by rounded distance; ties have no guaranteed order.
There is currently no pagination or selectable sort for this endpoint.

### List a restaurant's dishes

```bash
curl -i http://localhost:8080/api/restaurants/1/dishes
```

Returns dish objects in ID order, without pagination. An existing restaurant
without dishes returns `200` with `[]`. Invalid IDs return `400`; missing
restaurants return `404`.

## Error responses

| Status | Meaning |
| --- | --- |
| `400 Bad Request` | Invalid ID, query parameter, JSON body, or required field |
| `404 Not Found` | Requested dish/restaurant does not exist, or route is unknown |
| `413 Request Entity Too Large` | Request body exceeds the 64 KiB limit |
| `500 Internal Server Error` | Unexpected server/database failure |

Example validation response body (plain text):

```text
limit must be between 1 and 100
```

Database failures handled by `serverError` return `internal server error` to the
client; diagnostic details are logged on the server. Clients must check the HTTP
status before trying to decode the response as JSON.

## Tests

```bash
make test
make test-integration
```

The first command excludes PostgreSQL integration tests. The second requires
the configured, migrated, and seeded `cheap_eats_test` database. GitHub Actions
runs both categories, creating and seeding its own temporary database for the
integration job.
