# Kasoku Docker Compose Deployment

Production-ready Docker Compose setup for self-hosting Kasoku with automatic SSL, PostgreSQL database, and persistent storage.

## ✨ Features

- 🔒 **Automatic HTTPS** via Caddy with Let's Encrypt
- 🗄️ **PostgreSQL** database with health checks
- 📦 **Persistent volumes** for data and storage
- 🔄 **Automatic restarts** and health monitoring
- 🔐 **Security headers** and best practices
- 📊 **Built-in metrics** endpoint
- 💾 **Automated backups** with retention policy

## 🚀 Quick Start

### Prerequisites

- Linux server (Ubuntu, Debian, CentOS, etc.)
- Docker and Docker Compose installed
- Domain name pointing to your server
- At least 2GB RAM recommended

### One-Command Setup

```bash
sudo ./scripts/setup.sh
```

This interactive script will:
1. Install Docker (if needed)
2. Configure your domain and SSL
3. Generate secure secrets
4. Set up storage (local or S3)
5. Optionally configure Stripe billing

### Manual Setup

If you prefer to configure manually:

1. **Clone the repository**:
   ```bash
   git clone https://github.com/thebushidocollective/kasoku
   cd kasoku/deployment/docker-compose
   ```

2. **Configure environment**:
   ```bash
   cp .env.example .env
   nano .env  # Edit with your settings
   ```

   **Minimum required changes:**
   - `DOMAIN`: Your domain name
   - `DB_PASSWORD`: Generate with `openssl rand -base64 32`
   - `JWT_SECRET`: Generate with `openssl rand -hex 32`

3. **Point your domain to the server**:
   ```bash
   # Add A record:
   kasoku.example.com → YOUR_SERVER_IP
   ```

4. **Start Kasoku**:
   ```bash
   docker-compose -f docker-compose.production.yml up -d
   ```

5. **Check logs**:
   ```bash
   docker-compose -f docker-compose.production.yml logs -f
   ```

6. **Visit your installation**:
   ```
   https://your-domain.com
   ```

   SSL certificates are provisioned automatically on first access (may take 1-2 minutes).

## 📋 Configuration Options

### Storage Options

#### Local Storage (Default)
Best for self-hosted deployments:
```bash
STORAGE_TYPE=local
```

#### S3 Storage
For cloud/scalable deployments:
```bash
STORAGE_TYPE=s3
S3_BUCKET=your-bucket-name
S3_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-key
AWS_SECRET_ACCESS_KEY=your-secret
```

### Billing (Optional)

Leave empty for self-hosted without billing:
```bash
STRIPE_SECRET_KEY=
STRIPE_WEBHOOK_SECRET=
```

Or configure Stripe for paid plans:
```bash
STRIPE_SECRET_KEY=sk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_INDIVIDUAL_PRICE_ID=price_...
STRIPE_TEAM_PRICE_ID=price_...
```

## 🔧 Management

### Start/Stop

```bash
# Start
docker-compose -f docker-compose.production.yml up -d

# Stop
docker-compose -f docker-compose.production.yml down

# Restart
docker-compose -f docker-compose.production.yml restart

# View logs
docker-compose -f docker-compose.production.yml logs -f

# View specific service logs
docker-compose -f docker-compose.production.yml logs -f kasoku-server
```

### Updates

```bash
# Pull latest images
docker-compose -f docker-compose.production.yml pull

# Restart with new images
docker-compose -f docker-compose.production.yml up -d
```

### Backups

**Automated backups:**
```bash
# Run backup script
./scripts/backup.sh

# Schedule daily backups (cron)
0 2 * * * /path/to/kasoku/deployment/docker-compose/scripts/backup.sh
```

Backups are stored in `./backups/` and include:
- PostgreSQL database dump
- Storage volume (cache artifacts)
- Configuration (.env file)

**Restore from backup:**
```bash
# Restore database
gunzip < backups/postgres_TIMESTAMP.sql.gz | \
  docker-compose -f docker-compose.production.yml exec -T postgres psql -U kasoku kasoku

# Restore storage
docker run --rm -v $(pwd)_storage_data:/data -v ./backups:/backup \
  alpine tar xzf /backup/storage_TIMESTAMP.tar.gz -C /data
```

## 🔍 Monitoring

### Health Checks

All services have health checks configured:
```bash
docker-compose -f docker-compose.production.yml ps
```

### Service Status

```bash
# Check individual service
docker-compose -f docker-compose.production.yml exec kasoku-server wget -qO- http://localhost:8080/health

# Check Caddy metrics
curl http://localhost:2019/metrics
```

### Logs

```bash
# All services
docker-compose -f docker-compose.production.yml logs -f

# Specific service
docker-compose -f docker-compose.production.yml logs -f kasoku-server

# Tail last 100 lines
docker-compose -f docker-compose.production.yml logs --tail=100
```

## 🐛 Troubleshooting

### SSL Certificate Issues

**Symptom**: "certificate not found" or connection errors

**Solutions**:
1. Verify DNS is pointing to your server: `dig your-domain.com`
2. Check Caddy logs: `docker-compose logs caddy`
3. Ensure ports 80 and 443 are open in firewall
4. Wait 1-2 minutes for initial cert provisioning

### Database Connection Errors

**Symptom**: "failed to connect to database"

**Solutions**:
1. Check PostgreSQL is healthy: `docker-compose ps`
2. Verify DB_PASSWORD in .env matches
3. Check PostgreSQL logs: `docker-compose logs postgres`
4. Restart services: `docker-compose restart`

### Out of Disk Space

**Symptom**: "no space left on device"

**Solutions**:
```bash
# Check disk usage
df -h

# Clean Docker system
docker system prune -a

# Clean old backups
rm -f backups/postgres_$(date -d '30 days ago' +%Y%m%d)*.sql.gz

# Check volume sizes
docker system df -v
```

### Port Already in Use

**Symptom**: "address already in use"

**Solutions**:
```bash
# Find what's using the port
sudo lsof -i :80
sudo lsof -i :443

# Change ports in .env
HTTP_PORT=8080
HTTPS_PORT=8443

# Or stop the conflicting service
sudo systemctl stop apache2
```

## 🔐 Security Best Practices

### Before Going to Production

- [ ] Change all default passwords in `.env`
- [ ] Generate new `JWT_SECRET` and `DB_PASSWORD`
- [ ] Set up firewall (allow only 80, 443, and SSH)
- [ ] Enable automated backups
- [ ] Set up monitoring/alerts
- [ ] Review Caddy security headers
- [ ] Keep Docker and images updated

### Firewall Configuration

```bash
# UFW (Ubuntu/Debian)
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# Firewalld (CentOS/RHEL)
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### Regular Maintenance

```bash
# Update system packages
sudo apt update && sudo apt upgrade -y  # Ubuntu/Debian
sudo yum update -y                      # CentOS/RHEL

# Update Docker images
docker-compose -f docker-compose.production.yml pull
docker-compose -f docker-compose.production.yml up -d

# Clean up unused resources
docker system prune -f
```

## 📊 Resource Requirements

### Minimum
- **CPU**: 1 vCore
- **RAM**: 2 GB
- **Disk**: 20 GB
- **Network**: 100 Mbps

### Recommended
- **CPU**: 2 vCores
- **RAM**: 4 GB
- **Disk**: 50 GB SSD
- **Network**: 1 Gbps

### Scaling

For larger deployments:
- Use S3 storage instead of local
- Increase database resources (PostgreSQL memory)
- Scale kasoku-server horizontally (load balancer + multiple instances)
- Consider Kubernetes deployment (see `deployment/kubernetes/`)

## 🆘 Getting Help

- **Documentation**: https://github.com/thebushidocollective/kasoku
- **Issues**: https://github.com/thebushidocollective/kasoku/issues
- **Discussions**: https://github.com/thebushidocollective/kasoku/discussions

## 📜 License

Same as main Kasoku project.
