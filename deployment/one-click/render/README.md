# Deploy Kasoku to Render

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/thebushidocollective/brisk)

## One-Click Deployment

Click the button above to deploy Kasoku to Render with a single click.

## What Gets Deployed

The blueprint automatically creates:

- **Kasoku Server**: Go backend service with health checks
- **Kasoku Web**: Next.js frontend application
- **PostgreSQL Database**: Managed PostgreSQL instance
- **Persistent Disk**: 10GB disk for cache storage

## Manual Deployment

### Prerequisites

- Render account ([sign up here](https://render.com))
- GitHub repository with Kasoku code

### Step 1: Create Services

#### Create PostgreSQL Database

1. Go to Render Dashboard → New → PostgreSQL
2. Configure:
   - Name: `kasoku-db`
   - Database: `kasoku`
   - User: `kasoku`
   - Region: Choose closest to your users
   - Plan: Starter ($7/month) or higher

3. Note the Internal Database URL for later

#### Create Server Service

1. Go to Render Dashboard → New → Web Service
2. Connect your GitHub repository
3. Configure:
   - Name: `kasoku-server`
   - Environment: Docker
   - Region: Same as database
   - Branch: `main`
   - Dockerfile Path: `./server/Dockerfile`
   - Docker Context: `./server`

4. Set Environment Variables:
   ```
   PORT=8080
   DB_DRIVER=postgres
   DB_DSN=<internal-database-url>
   JWT_SECRET=<generate-secure-secret>
   STORAGE_TYPE=local
   STORAGE_PATH=/data
   ```

5. Add Persistent Disk:
   - Name: `kasoku-storage`
   - Mount Path: `/data`
   - Size: 10GB

6. Health Check Path: `/health`

#### Create Web Service

1. Go to Render Dashboard → New → Web Service
2. Connect your GitHub repository
3. Configure:
   - Name: `kasoku-web`
   - Environment: Docker
   - Region: Same as server
   - Branch: `main`
   - Dockerfile Path: `./web/Dockerfile`
   - Docker Context: `./web`

4. Set Environment Variables:
   ```
   NEXT_PUBLIC_API_URL=https://kasoku-server.onrender.com
   ```

### Step 2: Configure Blueprint (Alternative)

Instead of manual setup, use the `render.yaml` blueprint:

1. Fork the Kasoku repository
2. Add `render.yaml` to repository root
3. Connect to Render via Dashboard → New → Blueprint
4. Render will automatically create all services

### Step 3: Configure Custom Domain

1. Go to service settings
2. Navigate to "Custom Domains"
3. Add your domain (e.g., `kasoku.example.com`)
4. Configure DNS:
   ```
   CNAME kasoku.example.com -> <your-service>.onrender.com
   ```

## Architecture on Render

```
┌──────────────────────────────┐
│         Render Cloud         │
├──────────────────────────────┤
│                              │
│  ┌────────────────────────┐  │
│  │    kasoku-web          │  │
│  │    (Next.js)           │  │
│  │    Port: 3000          │  │
│  └──────────┬─────────────┘  │
│             │                │
│  ┌──────────▼─────────────┐  │
│  │    kasoku-server       │  │
│  │    (Go API)            │  │
│  │    Port: 8080          │  │
│  │    /health ✓           │  │
│  └──────────┬─────────────┘  │
│             │                │
│  ┌──────────▼─────────────┐  │
│  │    kasoku-db           │  │
│  │    (PostgreSQL)        │  │
│  │    Port: 5432          │  │
│  └────────────────────────┘  │
│                              │
│  ┌────────────────────────┐  │
│  │  Persistent Disk       │  │
│  │  /data (10GB)          │  │
│  └────────────────────────┘  │
└──────────────────────────────┘
```

## Storage Options

### Local Storage (Default)

Render provides persistent disks:
- Mounted at `/data`
- Survives deployments
- Backed up automatically
- Size: 10GB minimum

### S3 Storage (Recommended for Production)

For better scalability:

```bash
STORAGE_TYPE=s3
S3_BUCKET=kasoku-cache
S3_REGION=us-east-1
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=...
```

## Pricing Estimate

Render pricing (as of 2024):

### Starter Configuration
- Web Service (Server): $7/month
- Web Service (Web): $7/month
- PostgreSQL Starter: $7/month
- Persistent Disk (10GB): $0.25/GB = $2.50/month
- **Total**: ~$23.50/month

### Production Configuration
- Web Service (Standard): $25/month × 2 = $50
- PostgreSQL Standard: $20/month
- Persistent Disk (50GB): $12.50/month
- **Total**: ~$82.50/month

## Environment Variables

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_DRIVER` | Database driver | `postgres` |
| `DB_DSN` | Database connection string | From Render database |
| `JWT_SECRET` | JWT signing secret | Generate with `openssl rand -hex 32` |
| `STORAGE_TYPE` | Storage backend | `local` or `s3` |

### Optional Variables (Billing)

| Variable | Description |
|----------|-------------|
| `STRIPE_SECRET_KEY` | Stripe secret key |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook secret |
| `STRIPE_INDIVIDUAL_PRICE_ID` | Individual plan price ID |
| `STRIPE_TEAM_PRICE_ID` | Team plan price ID |
| `STRIPE_SUCCESS_URL` | Redirect after successful payment |
| `STRIPE_CANCEL_URL` | Redirect after cancelled payment |
| `STRIPE_RETURN_URL` | Return URL for billing portal |

## Monitoring and Logs

### Viewing Logs

1. Dashboard → Select Service → Logs tab
2. Real-time log streaming
3. Search and filter capabilities

### Metrics

Render provides built-in metrics:
- CPU usage
- Memory usage
- Request count
- Response times
- Error rates

Access via: Dashboard → Service → Metrics

### Alerts

Configure alerts for:
- Service health check failures
- High resource usage
- Error rate spikes

## Backup and Recovery

### Database Backups

Render automatically backs up PostgreSQL:
- Daily backups for Starter plan
- Point-in-time recovery for Standard/Pro plans
- Restore via Dashboard → Database → Backups

### Manual Database Backup

```bash
# Install Render CLI
brew install render

# Export database
render db:backup kasoku-db

# Or use pg_dump
pg_dump $DATABASE_URL > backup.sql
```

### Disk Backups

Persistent disks are backed up automatically:
- Daily snapshots
- Restore from Dashboard → Disk → Snapshots

## Auto-Deploy from Git

Render automatically deploys on git push:

1. Go to service settings
2. Enable "Auto-Deploy"
3. Select branch (e.g., `main`)
4. Push code triggers deployment

### Deploy Hooks

Trigger manual deploys via webhook:

```bash
curl -X POST https://api.render.com/deploy/srv-xxx?key=yyy
```

## Scaling

### Vertical Scaling

Upgrade service plan:
1. Dashboard → Service → Settings
2. Change instance type
3. Apply changes (causes restart)

### Horizontal Scaling

For multiple instances:
1. Requires Standard plan or higher
2. Dashboard → Service → Scaling
3. Increase instance count
4. Note: Use S3 storage for stateless operation

## Health Checks

Render monitors `/health` endpoint:
- Checks every 30 seconds
- Restarts unhealthy services
- Configurable timeout and retries

Configure in service settings:
```yaml
healthCheckPath: /health
```

## SSL/TLS

Render provides automatic SSL:
- Free SSL certificates
- Auto-renewal
- HTTP → HTTPS redirect
- Custom domain support

## Troubleshooting

### Build Failures

Check build logs:
1. Dashboard → Service → Events
2. Review error messages
3. Common issues:
   - Missing dependencies
   - Incorrect Dockerfile path
   - Build timeout (increase in settings)

### Database Connection Issues

Verify connection string:
```bash
# In service shell
echo $DATABASE_URL

# Test connection
psql $DATABASE_URL -c "SELECT version();"
```

### Service Won't Start

1. Check health check endpoint
2. Review environment variables
3. Verify port configuration (must bind to `0.0.0.0:$PORT`)
4. Check logs for errors

### Disk Space Issues

```bash
# SSH into service
render shell kasoku-server

# Check disk usage
df -h /data

# Clean up cache
rm -rf /data/cache/*
```

## Advanced Features

### Preview Environments

Automatic preview deployments for PRs:
1. Settings → Enable "Preview Environments"
2. Each PR gets unique URL
3. Separate database and storage

### Private Services

Keep services internal:
1. Service Settings → Networking
2. Select "Private"
3. Only accessible within Render network

### Background Workers

Create worker services:
```yaml
services:
  - type: worker
    name: kasoku-worker
    env: docker
    dockerfilePath: ./worker/Dockerfile
```

### Cron Jobs

Schedule tasks:
```yaml
services:
  - type: cron
    name: kasoku-cleanup
    schedule: "0 2 * * *"  # 2 AM daily
    dockerfilePath: ./cleanup/Dockerfile
```

## CI/CD Integration

### GitHub Actions

`.github/workflows/deploy.yml`:
```yaml
name: Deploy to Render
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Trigger Render Deploy
        env:
          RENDER_DEPLOY_HOOK: ${{ secrets.RENDER_DEPLOY_HOOK }}
        run: |
          curl -X POST $RENDER_DEPLOY_HOOK
```

## Security Best Practices

1. **Use Environment Variables**: Never commit secrets
2. **Enable 2FA**: On Render account
3. **Restrict Access**: Use private services when possible
4. **Regular Updates**: Keep dependencies updated
5. **Monitor Logs**: Set up alerts for suspicious activity
6. **Database Access**: Use internal URLs only
7. **SSL/TLS**: Always enabled by default

## Support

- Render Docs: https://render.com/docs
- Render Community: https://community.render.com
- Render Status: https://status.render.com
- Kasoku Issues: https://github.com/thebushidocollective/brisk/issues
