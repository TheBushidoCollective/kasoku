# Deploy Kasoku to Fly.io

Deploy Kasoku to Fly.io with global edge deployment and automatic scaling.

## Prerequisites

- Fly.io account ([sign up here](https://fly.io))
- Fly CLI installed: `curl -L https://fly.io/install.sh | sh`
- Docker installed (for local testing)

## Quick Start

```bash
# Login to Fly.io
fly auth login

# Deploy everything with one script
./deployment/one-click/fly/deploy.sh
```

## Manual Deployment

### Step 1: Create PostgreSQL Database

```bash
# Create Postgres cluster
fly postgres create \
  --name kasoku-db \
  --region sjc \
  --initial-cluster-size 1 \
  --vm-size shared-cpu-1x \
  --volume-size 10

# Note the connection string provided
```

### Step 2: Deploy Server

```bash
# Create server app
fly apps create kasoku-server

# Create storage volume
fly volumes create kasoku_storage \
  --region sjc \
  --size 10 \
  --app kasoku-server

# Set secrets
fly secrets set \
  JWT_SECRET=$(openssl rand -hex 32) \
  DB_DSN="postgres://..." \
  --app kasoku-server

# Optional: Stripe billing
fly secrets set \
  STRIPE_SECRET_KEY="sk_live_..." \
  STRIPE_WEBHOOK_SECRET="whsec_..." \
  STRIPE_INDIVIDUAL_PRICE_ID="price_..." \
  STRIPE_TEAM_PRICE_ID="price_..." \
  --app kasoku-server

# Deploy server
fly deploy \
  --config deployment/one-click/fly/fly.toml \
  --app kasoku-server
```

### Step 3: Deploy Web Dashboard

```bash
# Create web app
fly apps create kasoku-web

# Set environment variables
fly secrets set \
  NEXT_PUBLIC_API_URL="https://kasoku-server.fly.dev" \
  --app kasoku-web

# Deploy web
fly deploy \
  --config deployment/one-click/fly/fly-web.toml \
  --app kasoku-web
```

### Step 4: Configure Custom Domain

```bash
# Add certificate for custom domain
fly certs add kasoku.example.com --app kasoku-web

# Configure DNS
# Add these records to your DNS provider:
# A     kasoku.example.com    -> (Fly.io IPv4)
# AAAA  kasoku.example.com    -> (Fly.io IPv6)

# Verify
fly certs show kasoku.example.com --app kasoku-web
```

## Architecture on Fly.io

```
┌─────────────────────────────────────┐
│         Fly.io Global Network       │
├─────────────────────────────────────┤
│                                     │
│  ┌───────────────────────────────┐  │
│  │   kasoku-web (Next.js)        │  │
│  │   Region: sjc, ams, syd       │  │
│  │   Auto-scale: 1-10 machines   │  │
│  └──────────────┬────────────────┘  │
│                 │                   │
│  ┌──────────────▼────────────────┐  │
│  │   kasoku-server (Go API)      │  │
│  │   Region: sjc                 │  │
│  │   Health: /health             │  │
│  └──────────────┬────────────────┘  │
│                 │                   │
│  ┌──────────────▼────────────────┐  │
│  │   kasoku-db (PostgreSQL)      │  │
│  │   Region: sjc                 │  │
│  │   HA: Primary + Replica       │  │
│  └───────────────────────────────┘  │
│                                     │
│  ┌───────────────────────────────┐  │
│  │   Volume: kasoku_storage      │  │
│  │   Size: 10GB                  │  │
│  │   Mount: /data                │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

## Configuration Files

### Server Configuration (fly.toml)

```toml
app = "kasoku-server"
primary_region = "sjc"

[build]
  dockerfile = "server/Dockerfile"

[http_service]
  internal_port = 8080
  force_https = true

[mounts]
  source = "kasoku_storage"
  destination = "/data"
```

### Web Configuration (fly-web.toml)

```toml
app = "kasoku-web"
primary_region = "sjc"

[build]
  dockerfile = "web/Dockerfile"

[http_service]
  internal_port = 3000
  force_https = true
```

## Regions and Scaling

### Available Regions

Choose regions close to your users:
- `sjc` - San Jose, CA (US West)
- `iad` - Ashburn, VA (US East)
- `lhr` - London (Europe)
- `nrt` - Tokyo (Asia)
- `syd` - Sydney (Australia)

See all regions: `fly platform regions`

### Multi-Region Deployment

Deploy to multiple regions:

```bash
# Scale server to multiple regions
fly scale count 3 \
  --region sjc,iad,lhr \
  --app kasoku-server

# Scale web globally
fly scale count 5 \
  --region sjc,iad,lhr,nrt,syd \
  --app kasoku-web
```

### Auto-Scaling

Configure auto-scaling in `fly.toml`:

```toml
[http_service]
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 1
```

## Storage Options

### Fly Volumes (Default)

Persistent volumes:
- Local SSD storage
- Fast performance
- Regional (not replicated)
- Survives deployments

```bash
# Create volume
fly volumes create kasoku_storage \
  --size 10 \
  --region sjc

# List volumes
fly volumes list

# Extend volume
fly volumes extend <volume-id> --size 20
```

### S3 Storage (Recommended for Multi-Region)

For global deployments:

```bash
fly secrets set \
  STORAGE_TYPE=s3 \
  S3_BUCKET=kasoku-cache \
  S3_REGION=us-east-1 \
  AWS_ACCESS_KEY_ID=AKIA... \
  AWS_SECRET_ACCESS_KEY=... \
  --app kasoku-server
```

## Pricing Estimate

Fly.io pricing (as of 2024):

### Free Tier
- 3 shared-cpu-1x VMs (256MB RAM)
- 3GB persistent volume storage
- 160GB outbound data transfer

### Starter Configuration (~$10/month)
- Server: 1x shared-cpu-1x, 512MB RAM (~$3)
- Web: 1x shared-cpu-1x, 256MB RAM (~$2)
- PostgreSQL: Starter (~$0 on free tier or $15/month)
- Volume: 10GB (~$2)

### Production Configuration (~$50/month)
- Server: 2x shared-cpu-2x, 1GB RAM (~$20)
- Web: 3x shared-cpu-1x, 512MB RAM (~$15)
- PostgreSQL: Standard (~$30)
- Volume: 20GB (~$4)

Calculate your costs: https://fly.io/docs/about/pricing/

## Monitoring and Logs

### View Logs

```bash
# Tail logs
fly logs --app kasoku-server

# Last 100 lines
fly logs --app kasoku-server -n 100

# Filter by instance
fly logs --app kasoku-server --instance <instance-id>
```

### Metrics Dashboard

```bash
# Open metrics dashboard
fly dashboard --app kasoku-server

# Or visit: https://fly.io/apps/kasoku-server/metrics
```

### Monitoring Endpoints

Fly provides:
- Health checks
- Request metrics
- Response times
- Error rates
- VM metrics (CPU, memory, disk)

### Alerts

Set up alerts via Dashboard:
1. Go to app settings
2. Configure alerts for:
   - Health check failures
   - High resource usage
   - Certificate expiration

## Secrets Management

### Set Secrets

```bash
# Set individual secret
fly secrets set JWT_SECRET=xxx --app kasoku-server

# Set multiple secrets
fly secrets set \
  STRIPE_SECRET_KEY=sk_... \
  STRIPE_WEBHOOK_SECRET=whsec_... \
  --app kasoku-server

# Import from file
fly secrets import < secrets.env --app kasoku-server
```

### List Secrets

```bash
# List secret names (not values)
fly secrets list --app kasoku-server
```

### Remove Secrets

```bash
fly secrets unset STRIPE_SECRET_KEY --app kasoku-server
```

## Database Management

### Connect to Database

```bash
# Connect via psql
fly postgres connect --app kasoku-db

# Or proxy locally
fly proxy 5432 --app kasoku-db
psql postgres://kasoku:password@localhost:5432/kasoku
```

### Backup Database

```bash
# Create backup
fly postgres backup create --app kasoku-db

# List backups
fly postgres backup list --app kasoku-db

# Restore backup
fly postgres backup restore <backup-id> --app kasoku-db
```

### Database High Availability

```bash
# Add replica
fly postgres create \
  --name kasoku-db-replica \
  --region iad \
  --fork-from kasoku-db

# Scale replicas
fly postgres ha --app kasoku-db
```

## Deployment Strategies

### Blue-Green Deployment

```bash
# Deploy to canary
fly deploy --strategy canary

# Monitor and promote
fly releases list
fly releases promote <release-version>
```

### Rolling Deployment (Default)

```bash
# Standard deploy
fly deploy

# Deploy with specific strategy
fly deploy --strategy rolling
```

### Immediate Deployment

```bash
# Deploy all machines immediately
fly deploy --strategy immediate
```

## CI/CD Integration

### GitHub Actions

`.github/workflows/fly.yml`:
```yaml
name: Deploy to Fly.io

on:
  push:
    branches: [main]

jobs:
  deploy-server:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - uses: superfly/flyctl-actions/setup-flyctl@master

      - name: Deploy Server
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
        run: |
          flyctl deploy \
            --config deployment/one-click/fly/fly.toml \
            --app kasoku-server

  deploy-web:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - uses: superfly/flyctl-actions/setup-flyctl@master

      - name: Deploy Web
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
        run: |
          flyctl deploy \
            --config deployment/one-click/fly/fly-web.toml \
            --app kasoku-web
```

## Advanced Features

### Private Networking

Connect services privately:

```bash
# Services in same org can communicate via:
# http://kasoku-server.internal:8080

# Update web config
fly secrets set \
  NEXT_PUBLIC_API_URL="http://kasoku-server.internal:8080" \
  --app kasoku-web
```

### WireGuard VPN

Access your apps securely:

```bash
# Create WireGuard tunnel
fly wireguard create

# Connect to internal network
fly proxy 5432 --app kasoku-db
```

### Redis Cache

Add Redis for caching:

```bash
# Create Redis
fly redis create \
  --name kasoku-redis \
  --region sjc

# Get connection URL
fly redis status kasoku-redis

# Set in app
fly secrets set REDIS_URL=redis://... --app kasoku-server
```

## Troubleshooting

### App Won't Start

```bash
# Check status
fly status --app kasoku-server

# View logs
fly logs --app kasoku-server

# SSH into machine
fly ssh console --app kasoku-server

# Check environment
fly ssh console --app kasoku-server -C "env | grep -E '(DB|JWT|STORAGE)'"
```

### Database Connection Issues

```bash
# Test connection
fly ssh console --app kasoku-server -C "nc -zv kasoku-db.internal 5432"

# Check DNS
fly ssh console --app kasoku-server -C "nslookup kasoku-db.internal"

# Verify credentials
fly postgres list --app kasoku-db
```

### Storage Issues

```bash
# Check volume status
fly volumes list --app kasoku-server

# Check disk usage
fly ssh console --app kasoku-server -C "df -h"

# Mount issues - restart app
fly apps restart kasoku-server
```

### SSL Certificate Issues

```bash
# Check certificate
fly certs show kasoku.example.com --app kasoku-web

# Renew certificate
fly certs check kasoku.example.com --app kasoku-web
```

## Performance Optimization

### VM Sizing

```bash
# Upgrade VM
fly scale vm shared-cpu-2x --memory 1024 --app kasoku-server

# View current scale
fly scale show --app kasoku-server
```

### Connection Pooling

For database connections:

```bash
# Create PgBouncer
fly postgres create \
  --name kasoku-db-bouncer \
  --consul-url <consul-url>
```

### CDN Integration

Use Fly's built-in CDN:
- Automatic for static assets
- Edge caching
- HTTP/2 and HTTP/3

## Security

### Network Security

```bash
# Restrict database access
fly postgres config update \
  --allow-external-access=false \
  --app kasoku-db
```

### SSL/TLS

- Automatic SSL certificates
- Auto-renewal
- Force HTTPS (configured in fly.toml)
- TLS 1.2+ only

### Secrets

- Never commit secrets to git
- Use `fly secrets` for sensitive data
- Secrets are encrypted at rest

## Support

- Fly.io Docs: https://fly.io/docs
- Community: https://community.fly.io
- Status: https://status.flyio.net
- Kasoku Issues: https://github.com/thebushidocollective/brisk/issues
