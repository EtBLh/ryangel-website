# RyAngel Commerce API Specification

Version: 0.1.0 · Base URL (prod): `https://api.ryangel.com/api` · Base URL (staging): `https://staging-api.ryangel.com/api`

## 1. Conventions
- All endpoints accept and return `application/json` unless stated. Timestamps use ISO 8601 in UTC.
- Authentication uses opaque Bearer tokens (256-bit random strings). Tokens are hashed at rest and stored in `admin.token` / `client.token` with matching `token_expiry`.
- Pagination: `page` (default 1) and `page_size` (default 20, max 100). Responses wrap collections in `{ "data": [...], "meta": { "page": 1, ... } }`.
- Errors: `{ "error": { "code": "string", "message": "string", "details": {...} } }`. Common HTTP codes: 200, 201, 204, 400, 401, 403, 404, 409, 422, 500.
- All product media assets are stored as relative filesystem paths (e.g., `/media/products/faiachun-2025/main.jpg`) in the DB. The API server streams files back from `/media/...` and always returns absolute URLs in responses.

## 2. Authentication & Session Flows

### 2.1 Admin Authentication
| Method | Path | Description |
| --- | --- | --- |
| POST | `/admin/login` | Exchange username/email + password for opaque access token. |
| POST | `/admin/logout` | Invalidate the current token. |
| GET | `/admin/me` | Fetch admin profile and permissions. |

**POST /admin/login**
```json
{
	"username": "ops_lead",
	"password": "hunter2"
}
```
```json
{
	"token": "tn_f3791b1e19bd4d8f9ce35b5a8b3d08e7",
	"token_type": "Bearer",
	"expires_in": 3600,
	"admin": {
		"admin_id": 1,
		"username": "ops_lead",
		"email": "ops@example.com",
		"is_active": true,
		"last_login": "2025-12-11T10:10:00Z"
	}
}
```

### 2.2 Client Authentication
| Method | Path | Description |
| --- | --- | --- |
| POST | `/clients/register` | Register new client with phone number and send OTP. |
| POST | `/clients/login` | Send SMS OTP to existing phone number for authentication. |
| POST | `/clients/verify-otp` | Verify OTP code and obtain bearer token. |
| POST | `/clients/resend-otp` | Resend OTP if previous one expired. |
| GET | `/clients/me` | Get client profile (requires auth). |
| PATCH | `/clients/me` | Update client profile fields. |

Phone numbers must be unique and valid. OTP codes are 6-digit numbers valid for 5 minutes. Maximum 3 OTP requests per hour per phone number.

**POST /clients/register**
```json
{
	"phone": "+1234567890",
	"email": "user@example.com",
	"username": "johndoe"
}
```
```json
{
	"client": {
		"client_id": 123,
		"phone": "+1234567890",
		"email": "user@example.com",
		"username": "johndoe",
		"is_active": true,
		"created_at": "2025-12-13T10:00:00Z"
	},
	"message": "Client registered and OTP sent to +1234567890",
	"otp_expires_in": 300
}
```

**POST /clients/login**
```json
{
	"phone": "+1234567890"
}
```
```json
{
	"message": "OTP sent to +1234567890",
	"otp_expires_in": 300
}
```

**POST /clients/verify-otp**
```json
{
	"phone": "+1234567890",
	"otp": "123456"
}
```
```json
{
	"token": "tn_f3791b1e19bd4d8f9ce35b5a8b3d08e7",
	"token_type": "Bearer",
	"expires_in": 3600,
	"client": {
		"client_id": 123,
		"phone": "+1234567890",
		"email": "user@example.com",
		"username": "johndoe",
		"is_active": true
	}
}
```

**PATCH /clients/me**
```json
{
	"email": "newemail@example.com",
	"username": "newusername"
}
```
```json
{
	"client": {
		"client_id": 123,
		"phone": "+1234567890",
		"email": "newemail@example.com",
		"username": "newusername",
		"is_active": true,
		"updated_at": "2025-12-13T11:00:00Z"
	}
}
```

**POST /clients/resend-otp**
```json
{
	"phone": "+1234567890"
}
```
```json
{
	"message": "OTP resent to +1234567890",
	"otp_expires_in": 300
}
```

## 3. Client Addresses
Backed by `client_address` table with `UNIQUE (client_id) WHERE is_default`. Only one default shipping address per client.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/clients/me/addresses` | List addresses. Query `?is_default=true` for default. |
| POST | `/clients/me/addresses` | Create new address. Optional `is_default`; if true, resets others. |
| PATCH | `/clients/me/addresses/{address_id}` | Update address. |
| DELETE | `/clients/me/addresses/{address_id}` | Remove; cannot delete if referenced by open orders. |
| POST | `/clients/me/addresses/{address_id}/set-default` | Idempotent default setter. |

Address payload
```json
{
	"address_type": "home",
	"address_line1": "221B Baker St",
	"address_line2": null,
	"city": "London",
	"state": "London",
	"postal_code": "NW1",
	"country": "UK",
	"phone": "+44-20-0000",
	"is_default": true
}
```

## 4. Ebuy Store Pickup Locations
Backed by enriched `ebuy_store` table containing contact info, address fields, and GPS coordinates for local pickup locations.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/stores` | Public list of all available ebuy pickup locations. Returns enriched store data including GPS coordinates. |
| GET | `/stores/{store_id}` | Detail view for a specific store by string store_id. |

**GET /stores**
Returns a paginated list of all ebuy pickup locations with full details.

Response:
```json
{
  "data": [
    {
      "store_id": "H115A",
      "store_name": "水坑尾（竹園圍A櫃）",
      "type": "locker",
      "office_hours": "24hrs",
      "address": "澳門竹園圍斜巷7號東慶新村地下L2舖eHubs櫃H115A",
      "address_en": "SHOP L2, EDF. TONG HENG SAN CHUN, NO.7 TRAVESSA DO PADRE SOARES, MACAU",
      "latitude": 22.19431,
      "longitude": 113.54347
    },
    {
      "store_id": "H120A",
      "store_name": "黑沙環（廣華A櫃）",
      "type": "locker",
      "office_hours": "24hrs",
      "address": "澳門黑沙環新街445A號廣華新邨14座地下C舖eHubs櫃H120A",
      "address_en": "SHOP C, EDF. KUONG WA, NO.445A RUA NOVA DA AREIA PRETA, MACAU",
      "latitude": 22.207211,
      "longitude": 113.556491
    }
  ]
}
```

**GET /stores/{store_id}**
Returns detailed information for a specific store by its string store_id.

Response:
```json
{
  "store_id": "H115A",
  "store_name": "水坑尾（竹園圍A櫃）",
  "type": "locker",
  "office_hours": "24hrs",
  "address": "澳門竹園圍斜巷7號東慶新村地下L2舖eHubs櫃H115A",
  "address_en": "SHOP L2, EDF. TONG HENG SAN CHUN, NO.7 TRAVESSA DO PADRE SOARES, MACAU",
  "latitude": 22.19431,
  "longitude": 113.54347
}
```

Error responses:
- `404 Not Found`: Store with the specified store_id does not exist

## 5. Products
Product catalogue stored in `products`, linked to `categories` via `product_categories`, and to `product_images` (one-to-many) where each record stores the authoritative image path served by the API.

### 5.1 Catalogue
| Method | Path | Description |
| --- | --- | --- |
| GET | `/products` | Supports filters: `q`, `category_id`, `product_type`, `is_featured`, `price_min/max`, `sku`. Sort by `sort=price|-created_at`. |
| GET | `/products/{product_id}` | Returns product plus related categories, media, inventory. |

Sample list response
```json
{
	"data": [
		{
			"product_id": 10,
			"product_name": "FaiAchun 2025",
			"product_type": "faiachun",
			"price": "1299.00",
			"compare_at_price": "1499.00",
			"quantity": 45,
			"is_featured": true,
			"available_sizes": ["v-rect", "square", "fat-v-rect"],
			"images": [
				{
					"image_id": 501,
					"url": "https://api.ryangel.com/v1/media/products/faiachun-2025/main.jpg",
					"is_primary": true,
					"alt_text": "Front view",
					"size_type": "v-rect",
					"sort_order": 1
				},
				{
					"image_id": 502,
					"url": "https://api.ryangel.com/v1/media/products/faiachun-2025/square.jpg",
					"is_primary": false,
					"alt_text": "Square format",
					"size_type": "square",
					"sort_order": 2
				}
			],
			"tags": ["limited", "bundle"],
			"categories": [{"category_id": 3, "category_name": "Seasonal"}]
		}
	],
	"meta": {"page": 1, "page_size": 20, "total": 120}
}
```

### 5.2 Admin Product Management
| Method | Path | Notes |
| --- | --- | --- |
| POST | `/admin/products` | Create product; `sku` unique, `product_type` must match `product_type_enum`. |
| PATCH | `/admin/products/{product_id}` | Partial update, includes `quantity`, `is_active`, `tags`. Updates `updated_at` trigger. |
| POST | `/admin/products/{product_id}/images` | Upload or register new image path; API stores relative path in `product_images.image_path`. |
| POST | `/admin/products/{product_id}/categories` | Overwrite mapping list of `category_id`s. |

### 5.3 Product Images
`product_images` table stores metadata per asset.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/products/{product_id}/images` | Returns ordered list with `url`, `is_primary`, `alt_text`, `sort_order`. |

Sample response
```json
[
  {
    "url": "/media/products/faiachun-2025/main.jpg",
    "is_primary": true,
    "alt_text": "Front view",
    "size_type": "v-rect",
    "sort_order": 1
  },
  {
    "url": "/media/products/faiachun-2025/side.jpg", 
    "is_primary": false,
    "alt_text": "Side view",
    "size_type": "square",
    "sort_order": 2
  }
]
```
| POST | `/admin/products/{product_id}/images` | Accepts multipart upload or JSON `{ "image_path": "/media/products/...", "is_primary": true, "alt_text": "Front" }`. Server saves file (if uploaded), stores relative path, and returns rendered absolute URL. Only one `is_primary` is allowed; subsequent requests flip existing primary to false. |
| PATCH | `/admin/products/{product_id}/images/{image_id}` | Update metadata or reorder by changing `sort_order`. |
| DELETE | `/admin/products/{product_id}/images/{image_id}` | Remove asset and delete underlying file if no longer referenced. |

Responses expose the public URL computed as `base_url + image_path`, guaranteeing the file is ultimately served by the API server itself (no third-party CDN dependency by default).

## 6. Categories
- `GET /categories` returns flat list with parent references.
- `GET /categories/tree` returns nested tree.
- `POST /admin/categories` and `PATCH /admin/categories/{id}` manage taxonomy. Prevent cycles by validating requested `parent_category_id`.

## 7. Discounts & Promotions
Discount behavior mirrors `discounts`, `discount_products`, `discount_categories`.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/discounts` | Search by `discount_code`, `is_active`, `discount_type`, `applies_to`. |
| POST | `/admin/discounts` | Create. Validate `start_date <= end_date` and `usage_limit >= used_count`. |
| PATCH | `/admin/discounts/{id}` | Update or deactivate. |
| POST | `/admin/discounts/{id}/products` | Attach product IDs. |
| POST | `/admin/discounts/{id}/categories` | Attach category IDs. |
| POST | `/discounts/apply` | Client applies code to current cart; response includes recalculated totals and `discount_id`. |

`discount_type = 'bxgy'` requires `buy_quantity`, `get_quantity`, and optionally `free_product_id` (defaults to the item being purchased when `applies_to_same_product = true`).

## 8. Cart & Checkout
Carts can be anonymous or associated with a client. Anonymous carts are identified by `cart_id` passed in `X-Cart-ID` header. Logged-in clients have their cart associated via `client_id`. Inactive anonymous carts (no updates for 30 days) should be cleared periodically.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/cart` | Returns current cart and items with calculated totals (subtotal, discount, total). Requires auth or `X-Cart-ID` header. |

Sample response
```json
{
  "items": [
    {
      "product_id": 1,
      "size_type": "v-rect",
      "quantity": 2,
      "added_at": "2025-12-14T22:44:27.233715Z",
      "product_name": "FaiAchun 2025",
      "unit_price": 1299.00,
      "stock_quantity": 45
    }
  ],
  "subtotal": 2598.00,
  "discount": 0.00,
  "total": 2598.00
}
```
| POST | `/cart/items` | Body `{ "product_id": 5, "size_type": "v-rect", "quantity": 2 }`. `size_type` optional, defaults to null. Adds/updates item. Requires auth or `X-Cart-ID`. |
| PATCH | `/cart/items/{cart_item_id}` | Adjust quantity. Reject `quantity < 1`. |
| DELETE | `/cart/items/{cart_item_id}` | Remove item. |
| POST | `/cart/apply-discount` | `{ "discount_code": "SPRING25" }`. Requires auth. |
| DELETE | `/cart/discount` | Removes applied discount. Requires auth. |
| POST | `/cart/checkout` | Creates order. Requires auth. |

Checkout request (requires Bearer token)
```json
{
	"shipping_address_id": 12,
	"payment_method": "mpay",
	"notes": "Leave with concierge",
	"ebuy_store_id": null
}
```

Alternative checkout for local pickup:
```json
{
	"ebuy_store_id": "H115A",
	"payment_method": "mpay",
	"notes": "Pickup from 水坑尾 store"
}
```
Server must enforce the DB constraint: either `shipping_address_id` **xor** `ebuy_store_id`. The `ebuy_store_id` is a string identifier referencing an ebuy pickup location.

**Manual payment expectation**
1. Checkout creates an order with `payment_status = "pending"` and stores the chosen `payment_method`.
2. Client pays externally via `mpay`, `boc`, or `bank_transfer` and records the transaction reference.
3. Client uploads payment proof (Section 10) for staff review.
4. Admin approves the proof, flipping `orders.payment_status` to `paid` and unlocking fulfilment.

## 9. Orders

### 9.1 Client-Facing
| Method | Path | Description |
| --- | --- | --- |
| GET | `/orders` | Lists client orders. Filter `order_status`, `date_from/to`. |
| GET | `/orders/{order_id}` | Includes nested `order_items`, `shipping_address` snapshot, payment status. |

### 9.2 Admin-Facing
| Method | Path | Description |
| --- | --- | --- |
| GET | `/admin/orders` | Global list with filters on status, store, payment method, `client_id`. |
| PATCH | `/admin/orders/{order_id}/status` | Body `{ "order_status": "shipped", "tracking_number": "SF123" }`. Allowed transitions follow business matrix (pending→confirmed→processing→shipped→delivered; cancelled/refunded terminal). |
| POST | `/admin/orders/{order_id}/refund` | Records refund, updates `payment_status`, writes `admin_notes`. |

### 9.3 Order Items Snapshot
- `order_items` capture `product_name`, `product_type`, `product_sku`, and `size_type` at purchase time to handle future catalog changes.
- Totals recomputed server-side: `total_price = (unit_price - discount_amount) * quantity`.
- `payment_reference` on `orders` stores the transaction identifier recorded by staff during proof approval (e.g., MPay receipt ID).

## 10. Payments (Manual Proof Flow)
No payment processor is integrated. Clients settle invoices externally and submit evidence through `payment_proofs`.

### 10.1 Client Actions
| Method | Path | Description |
| --- | --- | --- |
| POST | `/payments/proofs` | Authenticated clients upload payment evidence. Accepts multipart form (`file`, `order_id`, `payment_method`, `amount`, `transaction_reference`, optional `notes`). Server stores the file under `/media/payments/...`, records the relative `proof_path`, and responds with `{ "proof_id": 42, "status": "submitted" }`. Multiple submissions per order are allowed; the latest supersedes earlier attempts. |
| GET | `/payments/proofs/{proof_id}` | Client fetches their own submission metadata and an absolute download URL (computed as `base_url + proof_path`). |

### 10.2 Admin Actions
| Method | Path | Description |
| --- | --- | --- |
| GET | `/admin/payment-proofs` | Paginated listing filtered by `status`, `payment_method`, `order_id`, or `client_id`. Includes order totals to aid reconciliation. |
| PATCH | `/admin/payment-proofs/{proof_id}` | Approve or reject a proof. Body example: `{ "status": "approved", "review_notes": "Funds received", "payment_reference": "MPAY-8842" }`. Approval sets `payment_proofs.status = 'approved'`, stamps `reviewed_by` / `reviewed_at`, updates the linked order's `payment_status` to `paid`, and copies the provided `payment_reference` into `orders.payment_reference`. Rejection records the reason and leaves the order pending. |

`payment_proof_status_enum` defines the lifecycle: `submitted → approved|rejected`. Proof assets are always served by the API server, aligning with the `payment_proofs.proof_path` stored in the database.

## 11. Reporting & Analytics (Phase 2)
- `GET /admin/reports/sales?group_by=day` aggregates from `orders` and `order_items`.
- `GET /admin/reports/inventory` shows low-stock products (e.g., `quantity < safety_threshold`).
- `GET /admin/reports/discount-usage` summarises `used_count` per discount.

## 12. Webhooks (Optional)
- `order.created`: emitted after checkout with payload containing `order_id`, `client_id`, `total_amount`, `payment_status`.
- `order.updated`: triggers on status change.
- `inventory.low`: triggered when quantity crosses defined threshold. Webhook deliveries signed with shared secret.

## 13. Validation Rules Summary
- `orders`: request must respect shipping XOR pickup rule to satisfy DB constraint.
- `client_address`: API must demote previous default before setting new default.
- `discounts`: enforce date window and usage counters before allowing apply; increment `used_count` atomically with order creation.
- `cart`: enforce unique `(cart_id, product_id, size_type)` by using UPSERT semantics on add/update.
- `products`: prevent negative inventory; reject `quantity` updates that would break outstanding reservations (future enhancement: reserve table).

## 14. Error Catalogue
| Code | HTTP | Message | Notes |
| --- | --- | --- | --- |
| `AUTH_INVALID_CREDENTIALS` | 401 | Invalid username/email or password. | Applies to admins & clients. |
| `AUTH_TOKEN_EXPIRED` | 401 | Token expired. Refresh login. | |
| `ADDRESS_DEFAULT_CONFLICT` | 409 | Client already has default address. | Triggered when attempting to set second default without demoting first. |
| `ORDER_SHIPPING_CONFLICT` | 422 | Provide either shipping_address_id or ebuy_store_id. | Mirrors DB check constraint. |
| `DISCOUNT_NOT_APPLICABLE` | 422 | Discount cannot be applied. | Provide `details.reason`. |
| `INVENTORY_INSUFFICIENT` | 409 | Requested quantity exceeds stock. | Returned from cart add/update and checkout. |

## 15. Security & Observability
- Rate limit public endpoints to 60 req/min per IP; admin endpoints to 30 req/min.
- Enable audit logging for admin mutations (`admin_id`, action, payload hash, timestamp).
- All write endpoints require idempotency keys via `Idempotency-Key` header to prevent double-processing (e.g., cart checkout).
- Emit structured logs (JSON) with `trace_id` header propagation to support distributed tracing.

This specification stays aligned with the PostgreSQL schema in `db/000_init.sql`. Any schema change (new enum, column, constraint) must be reflected here before deployment.
