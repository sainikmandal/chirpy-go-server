# Chirpy Go Backend Server

## Introduction

**Chirpy** is a modern, secure, and extensible backend server for a microblogging platform, built with Go. It provides a robust RESTful API for user management, authentication (with JWT and refresh tokens), posting and retrieving chirps (short messages), and integration with external services like Polka for premium account upgrades.

Key features include:
- User registration, login, and secure password hashing
- JWT-based authentication and refresh token system
- Chirp creation, retrieval, filtering, and sorting
- Chirpy Red premium membership via secure webhook integration
- Admin and health endpoints for monitoring and maintenance
- Built-in metrics and reset functionality

The project is designed for extensibility, security, and clarity, making it a great foundation for learning Go web development or building your own social platform.

---

## API Endpoints

### `GET /api/chirps`

Returns a list of chirps.

#### Query Parameters

- `author_id` (optional): Filter chirps by the author's UUID.
- `sort` (optional): Sort the chirps by creation time.  
  - `asc` (default): Ascending order (oldest first)
  - `desc`: Descending order (newest first)

#### Examples

- Get all chirps (ascending order):
  ```
  GET /api/chirps
  ```
- Get all chirps (descending order):
  ```
  GET /api/chirps?sort=desc
  ```
- Get chirps by a specific author:
  ```
  GET /api/chirps?author_id=3311741c-680c-4546-99f3-fc9efac2036c
  ```
- Get chirps by a specific author, newest first:
  ```
  GET /api/chirps?author_id=3311741c-680c-4546-99f3-fc9efac2036c&sort=desc
  ```

---

### `GET /api/chirps/{chirpID}`

Returns a single chirp by its ID.

#### Example

```
GET /api/chirps/123e4567-e89b-12d3-a456-426614174000
```

---

### `POST /api/chirps`

Create a new chirp.  
**Requires JWT authentication.**

#### Request Body

```json
{
  "body": "Hello, Chirpy!"
}
```

#### Response

Returns the created chirp object.

---

### `DELETE /api/chirps/{chirpID}`

Delete a chirp by its ID.  
**Requires JWT authentication.**

---

### `POST /api/users`

Register a new user.

#### Request Body

```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

#### Response

Returns the created user object (without password).

---

### `PUT /api/users`

Update the authenticated user's email and/or password.  
**Requires JWT authentication.**

#### Request Body

```json
{
  "email": "newemail@example.com",
  "password": "newpassword"
}
```

---

### `POST /api/login`

Authenticate a user and receive a JWT and refresh token.

#### Request Body

```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

#### Response

```json
{
  "id": "user-uuid",
  "created_at": "...",
  "updated_at": "...",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "JWT_TOKEN",
  "refresh_token": "REFRESH_TOKEN"
}
```

---

### `POST /api/refresh`

Obtain a new JWT using a valid refresh token.

#### Headers

```
Authorization: Bearer REFRESH_TOKEN
```

---

### `POST /api/revoke`

Revoke a refresh token.

#### Headers

```
Authorization: Bearer REFRESH_TOKEN
```

---

### `POST /api/polka/webhooks`

Webhook endpoint for Polka to upgrade users to Chirpy Red.

#### Request Body

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
  }
}
```

#### Security

- Requires an API key in the `Authorization` header:
  ```
  Authorization: ApiKey YOUR_POLKA_KEY
  ```
- If the API key is missing or invalid, the endpoint responds with `401 Unauthorized`.

#### Behavior

- If the event is not `user.upgraded`, responds with `204 No Content`.
- If the user is upgraded successfully, responds with `204 No Content`.
- If the user is not found, responds with `404 Not Found`.

---

### `GET /api/healthz`

Health check endpoint.

---

### `GET /admin/metrics`

Returns HTML metrics for admin monitoring.

---

### `POST /admin/reset`

Resets the server's metrics counter.

---

## User Resource

All endpoints that return user resources now include the `is_chirpy_red` field:

```json
{
  "id": "user-uuid",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

---

## Environment Variables

- `DB_URL`: PostgreSQL connection string
- `JWT_SECRET`: ECDSA private key in PEM format for JWT signing
- `POLKA_KEY`: API key for Polka webhook authentication
- `PLATFORM`: (optional) Platform identifier 