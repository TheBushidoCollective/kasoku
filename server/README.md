# Kasoku Server

Self-hosted cache server for Kasoku. This is the same server that powers the hosted Kasoku service, minus the billing integration.

## Features

- Remote cache storage (local filesystem or S3)
- User authentication (JWT and API tokens)
- FIFO cache eviction with storage quotas
- Analytics tracking (cache hits, misses, time saved)
- Multi-user and team support
- RESTful API

## Quick Start with Docker

The easiest way to run Kasoku server is with Docker:

```bash
# Clone and navigate to server directory
cd server

# Start with Docker Compose (SQLite)
docker-compose up -d

# Or with PostgreSQL (edit docker-compose.yml to uncomment postgres sections)
docker-compose up -d kasoku-server-postgres postgres
```

The server will be available at `http://localhost:8080`

## Manual Setup

### Prerequisites

- Go 1.25 or later
- PostgreSQL or SQLite

### Installation

```bash
# Install dependencies
go mod download

# Copy and configure environment
cp .env.example .env
# Edit .env with your settings

# Build
go build -o kasoku-server ./cmd/kasoku-server

# Run
./kasoku-server
```

## Configuration

Configuration is done via environment variables. See `.env.example` for all options.

### Essential Configuration

```bash
# Server port
PORT=8080

# Database (sqlite or postgres)
DB_DRIVER=sqlite
DB_DSN=kasoku.db

# JWT secret (change in production!)
JWT_SECRET=your-very-long-random-secret-string

# Storage backend (local or s3)
STORAGE_TYPE=local
STORAGE_PATH=./storage
```

### Using S3 Storage

```bash
STORAGE_TYPE=s3
S3_BUCKET=your-bucket-name
S3_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

### Using PostgreSQL

```bash
DB_DRIVER=postgres
DB_DSN=host=localhost user=kasoku password=yourpassword dbname=kasoku port=5432 sslmode=disable
```

## API Endpoints

### Public Endpoints

- `POST /auth/register` - Create a new account
- `POST /auth/login` - Login and get JWT token
- `GET /health` - Health check

### Protected Endpoints

All protected endpoints require `Authorization: Bearer <token>` header.

#### Authentication

- `GET /auth/me` - Get current user info
- `POST /auth/tokens` - Create API token for CLI
- `GET /auth/tokens` - List API tokens
- `DELETE /auth/tokens/:id` - Revoke API token

#### Cache Operations

- `GET /cache` - List cache entries
- `PUT /cache/:hash` - Upload cache artifact
- `GET /cache/:hash` - Download cache artifact
- `DELETE /cache/:hash` - Delete cache entry

#### Analytics

- `GET /analytics` - Get analytics summary

## Usage

### 1. Create an Account

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "Your Name",
    "password": "your-secure-password"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your-secure-password"
  }'
```

This returns a JWT token.

### 3. Create API Token for CLI

```bash
curl -X POST http://localhost:8080/auth/tokens \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My CLI Token"
  }'
```

Save the returned token - it's only shown once!

### 4. Configure CLI

Configure your Kasoku CLI to use the self-hosted server:

```yaml
# kasoku.yaml
cache:
  remote:
    enabled: true
    type: "custom"
    url: "http://localhost:8080"
    token: "your-api-token"
```

## Storage Quotas

Default quotas by plan:

- **free**: 0 (no remote cache)
- **individual**: 2GB
- **team**: 50GB (configurable)

To change a user's plan, update the database directly:

```sql
UPDATE users SET plan = 'individual' WHERE email = 'user@example.com';
```

## Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

### View Logs

```bash
# Docker
docker-compose logs -f kasoku-server

# Manual
# Logs are written to stdout
```

## Security Considerations

1. **Change JWT_SECRET** - Use a long random string in production
2. **Use HTTPS** - Put the server behind a reverse proxy (nginx, Caddy, Traefik)
3. **Secure Database** - Use strong passwords and restrict access
4. **Backup** - Regular backups of database and storage
5. **Rate Limiting** - Consider adding rate limiting for API endpoints

## Differences from Hosted Service

The self-hosted server is identical to the hosted service except:

- No billing integration (Stripe)
- No usage-based pricing
- No automatic plan enforcement
- Manual user management

## Troubleshooting

### Database Connection Issues

- Check DB_DRIVER and DB_DSN environment variables
- Ensure PostgreSQL is running and accessible
- Check database credentials

### Storage Issues

- Ensure STORAGE_PATH directory exists and is writable
- For S3, verify AWS credentials and bucket permissions

### Authentication Errors

- Verify JWT_SECRET is set and consistent
- Check token expiration
- Ensure Authorization header is properly formatted

## License

Same as main Kasoku project.
