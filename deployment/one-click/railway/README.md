# Deploy Kasoku to Railway

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/kasoku)

## One-Click Deployment

Click the button above to deploy Kasoku to Railway with a single click.

## What Gets Deployed

The template will automatically create:

- **Kasoku Server**: Go backend service
- **Kasoku Web**: Next.js frontend
- **PostgreSQL Database**: Managed PostgreSQL instance
- **Redis** (optional): For caching and session management

## Manual Deployment

### Prerequisites

- Railway account ([sign up here](https://railway.app))
- Railway CLI (optional): `npm install -g @railway/cli`

### Step 1: Create a New Project

```bash
# Using Railway CLI
railway init

# Or create via web dashboard
```

### Step 2: Add PostgreSQL Database

```bash
railway add --database postgresql
```

### Step 3: Configure Environment Variables

Set the following environment variables in your Railway project:

#### Server Service

```bash
# Database (automatically set by Railway)
DATABASE_URL=${{Postgres.DATABASE_URL}}
DB_DRIVER=postgres

# Generate a secure JWT secret
JWT_SECRET=<generate-with-openssl-rand-hex-32>

# Storage configuration
STORAGE_TYPE=local
STORAGE_PATH=/data

# Optional: Stripe billing
STRIPE_SECRET_KEY=sk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_INDIVIDUAL_PRICE_ID=price_...
STRIPE_TEAM_PRICE_ID=price_...
STRIPE_SUCCESS_URL=https://your-domain.railway.app/dashboard?success=true
STRIPE_CANCEL_URL=https://your-domain.railway.app/pricing
STRIPE_RETURN_URL=https://your-domain.railway.app/dashboard
```

#### Web Service

```bash
NEXT_PUBLIC_API_URL=${{Server.RAILWAY_PUBLIC_DOMAIN}}
```

### Step 4: Deploy Services

```bash
# Deploy server
cd server
railway up

# Deploy web
cd ../web
railway up
```

### Step 5: Configure Custom Domain (Optional)

1. Go to your service settings in Railway dashboard
2. Click "Settings" → "Domains"
3. Add your custom domain
4. Configure DNS records as instructed

## Architecture on Railway

```
┌─────────────────┐
│   Railway       │
├─────────────────┤
│                 │
│  ┌──────────┐   │
│  │   Web    │   │  (Next.js)
│  │  :3000   │   │
│  └────┬─────┘   │
│       │         │
│  ┌────▼─────┐   │
│  │  Server  │   │  (Go API)
│  │  :8080   │   │
│  └────┬─────┘   │
│       │         │
│  ┌────▼─────┐   │
│  │PostgreSQL│   │
│  │  :5432   │   │
│  └──────────┘   │
│                 │
│  ┌──────────┐   │
│  │ Storage  │   │  (Volume)
│  │  /data   │   │
│  └──────────┘   │
└─────────────────┘
```

## Storage Considerations

Railway provides persistent storage through volumes:

1. Enable volumes in your service settings
2. Mount path: `/data`
3. Recommended size: 10GB+ for cache storage

**For production with high traffic**, consider using S3-compatible storage:

```bash
STORAGE_TYPE=s3
S3_BUCKET=kasoku-cache
S3_REGION=us-east-1
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
```

## Pricing Estimate

Railway pricing (as of 2024):

- **Hobby Plan**: $5/month + usage
  - $0.000463/GB-hour for memory
  - $0.000231/vCPU-hour for CPU
  - Database included

- **Estimated Monthly Cost**:
  - Small deployment (1-10 users): ~$10-20/month
  - Medium deployment (10-100 users): ~$30-50/month
  - Large deployment: Scale as needed

## Environment-Specific Configuration

### Development

```bash
railway environment development
railway variables set LOG_LEVEL=debug
```

### Production

```bash
railway environment production
railway variables set LOG_LEVEL=info
railway variables set DB_POOL_SIZE=20
```

## Monitoring and Logs

### View Logs

```bash
# CLI
railway logs

# Or view in dashboard
```

### Metrics

Railway provides built-in metrics:
- CPU usage
- Memory usage
- Network traffic
- Request count

Access via: Dashboard → Service → Metrics

## Backup and Recovery

### Database Backups

Railway automatically backs up PostgreSQL databases:
- Point-in-time recovery
- Access via: Dashboard → Database → Backups

### Manual Backup

```bash
# Export database
railway run pg_dump $DATABASE_URL > backup.sql

# Restore database
railway run psql $DATABASE_URL < backup.sql
```

### Volume Backups

Create manual backups of persistent volumes:

```bash
# SSH into service
railway run bash

# Create backup
tar czf /tmp/data-backup.tar.gz /data

# Download backup (use railway cli)
railway run cat /tmp/data-backup.tar.gz > data-backup.tar.gz
```

## Troubleshooting

### Service Won't Start

1. Check logs: `railway logs`
2. Verify environment variables are set
3. Ensure database is running
4. Check build logs for errors

### Database Connection Issues

```bash
# Test connection
railway run psql $DATABASE_URL -c "SELECT version();"

# Check DATABASE_URL is set
railway variables
```

### Out of Memory

Increase service memory:
1. Dashboard → Service → Settings
2. Adjust memory allocation
3. Redeploy service

### Storage Full

1. Check volume usage: `railway run df -h`
2. Clean up old cache: `railway run rm -rf /data/cache/*`
3. Increase volume size in settings

## Advanced Configuration

### Multiple Environments

```bash
# Create staging environment
railway environment create staging

# Deploy to staging
railway environment staging
railway up
```

### CI/CD Integration

Railway automatically deploys from Git:

1. Connect GitHub repository
2. Configure branch deployments
3. Set up PR previews (optional)

`.github/workflows/railway.yml`:
```yaml
name: Deploy to Railway
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Deploy to Railway
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
        run: |
          npm install -g @railway/cli
          railway up --service server
          railway up --service web
```

### Health Checks

Railway automatically monitors `/health` endpoint:
- Restart unhealthy services
- Configure in `railway.toml`

### Scaling

For horizontal scaling:
1. Dashboard → Service → Settings
2. Increase replica count
3. Note: Requires stateless architecture or external storage (S3)

## Support

- Railway Docs: https://docs.railway.app
- Railway Discord: https://discord.gg/railway
- Kasoku Issues: https://github.com/thebushidocollective/brisk/issues
