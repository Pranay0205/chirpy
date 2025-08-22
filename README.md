# Chirpy 🐦

A Twitter-like social media API built with Go, featuring user authentication, chirp management, and premium membership.

## Features

- **User Management**: Create accounts, login/logout with JWT authentication
- **Chirp System**: Post, read, update, and delete short messages (140 characters max)
- **Premium Membership**: Chirpy Red subscription with enhanced features
- **Webhook Integration**: Payment provider webhook handling
- **RESTful API**: Clean HTTP endpoints with JSON responses

## Tech Stack

- Go 1.24+
- PostgreSQL with sqlc
- JWT authentication
- Goose migrations

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL

### Installation

1. Clone and install dependencies:

```bash
git clone <repository-url>
cd chirpy
go mod download
```

2. Set up environment variables:

```bash
# Create .env file
DB_URL=postgres://username:password@localhost/chirpy?sslmode=disable
SECRET=your-jwt-secret-key
PLATFORM=dev
```

3. Run migrations and start:

```bash
cd sql/schema
goose postgres $DB_URL up
cd ../..
sqlc generate
make run
```

## Project Structure

```
chirpy/
├── internal/
│   ├── auth/           # Authentication utilities
│   └── database/       # Generated sqlc code
├── sql/
│   ├── schema/         # Database migrations
│   └── queries/        # SQL queries for sqlc
├── main.go             # Application entry point
├── handle_*.go         # HTTP handlers
├── api_types.go        # Request/Response structs
├── json.go             # JSON utilities
└── sqlc.yaml           # sqlc configuration
```

## API Endpoints

### Users

- `POST /api/users` - Create account
- `POST /api/login` - User login
- `PUT /api/users` - Update profile
- `POST /api/refresh` - Refresh token
- `POST /api/revoke` - Revoke token

### Chirps

- `GET /api/chirps` - Get chirps (supports `?sort=asc|desc` and `?author_id=uuid`)
- `GET /api/chirps/{id}` - Get specific chirp
- `POST /api/chirps` - Create chirp (auth required)
- `DELETE /api/chirps/{id}` - Delete own chirp (auth required)

### Admin

- `GET /api/healthz` - Health check
- `POST /api/polka/webhooks` - Payment webhooks
