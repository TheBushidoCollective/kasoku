# Kasoku Deployment Options

This directory contains multiple deployment strategies for Kasoku, optimized for different use cases and infrastructure requirements.

## 📦 Available Deployment Methods

### 1. Docker Compose (Recommended for Self-Hosted)

**Best for:** Personal use, small teams, single-server deployments

**Location:** `docker-compose/`

Simple, production-ready deployment with all services in one stack:
- ✅ One-command setup script
- ✅ Automatic SSL with Caddy
- ✅ Built-in PostgreSQL
- ✅ Local or S3 storage
- ✅ Automated backups

**Quick Start:**
```bash
cd deployment/docker-compose
./scripts/setup.sh
```

**Documentation:** [Docker Compose README](docker-compose/README.md)

---

### 2. Kubernetes with Helm (Recommended for Production)

**Best for:** Production deployments, high availability, enterprise

**Location:** `kubernetes/kasoku/`

Production-ready Helm chart with:
- ✅ High availability (multi-replica)
- ✅ Autoscaling (HPA)
- ✅ Built-in or external database
- ✅ S3 or local storage
- ✅ Ingress with SSL
- ✅ Pod disruption budgets
- ✅ Prometheus monitoring

**Quick Start:**
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install kasoku ./deployment/kubernetes/kasoku \
  --values deployment/kubernetes/kasoku/values-selfhosted.yaml \
  --set global.domain=kasoku.example.com
```

**Documentation:** [Helm Chart README](kubernetes/kasoku/README.md)

---

### 3. One-Click Cloud Deployments

**Best for:** Quick demos, managed infrastructure, global deployment

#### Railway

**Location:** `one-click/railway/`

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/kasoku)

- Automatic builds and deploys
- Managed PostgreSQL
- Built-in SSL
- ~$10-20/month

**Documentation:** [Railway README](one-click/railway/README.md)

#### Render

**Location:** `one-click/render/`

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/thebushidocollective/brisk)

- Blueprint-based deployment
- Auto-scaling
- Managed database
- ~$23/month starter

**Documentation:** [Render README](one-click/render/README.md)

#### Fly.io

**Location:** `one-click/fly/`

Global edge deployment with:
- Multi-region support
- Auto-scaling
- WireGuard VPN
- ~$10/month (includes free tier)

**Quick Start:**
```bash
./deployment/one-click/fly/deploy.sh
```

**Documentation:** [Fly.io README](one-click/fly/README.md)

---

## 🎯 Choosing the Right Deployment

### Use Docker Compose if you:
- Want simple, single-server deployment
- Are self-hosting on your own hardware
- Need quick local development environment
- Have basic scaling requirements (1-10 users)
- Prefer manual control over infrastructure

### Use Kubernetes/Helm if you:
- Need production-grade high availability
- Require horizontal scaling
- Want infrastructure as code
- Need multi-region deployment
- Have DevOps/SRE team
- Require compliance and audit logging

### Use One-Click Platforms if you:
- Want zero infrastructure management
- Need fast prototyping/demos
- Prefer managed services
- Want global CDN and edge deployment
- Need automatic SSL and domains
- Have budget for managed hosting

---

## 💰 Cost Comparison

| Deployment Method | Estimated Monthly Cost | Best For |
|-------------------|------------------------|----------|
| **Docker Compose** | $5-20 (VPS) | Self-hosted, single server |
| **Kubernetes** | $50-200+ (cluster) | Production, high traffic |
| **Railway** | $10-50 | Quick start, managed |
| **Render** | $23-80 | Startups, auto-scaling |
| **Fly.io** | $10-50 | Global edge, multi-region |

*Costs vary based on usage, storage, and traffic*

---

## 🏗️ Architecture Overview

All deployment methods include:

```
┌─────────────────────────────────────┐
│          Load Balancer / CDN        │
│         (SSL/TLS Termination)       │
└────────────────┬────────────────────┘
                 │
        ┌────────┴────────┐
        │                 │
┌───────▼──────┐  ┌──────▼───────┐
│   Web (UI)   │  │ Server (API) │
│   Next.js    │  │      Go      │
│   Port 3000  │  │   Port 8080  │
└──────────────┘  └──────┬───────┘
                         │
                  ┌──────┴──────┐
                  │             │
          ┌───────▼──────┐  ┌──▼────────┐
          │  PostgreSQL  │  │  Storage  │
          │   Database   │  │ (S3/Local)│
          └──────────────┘  └───────────┘
```

---

## 🚀 Quick Start Guide

### First Time Setup

1. **Choose your deployment method** based on the criteria above
2. **Clone the repository:**
   ```bash
   git clone https://github.com/thebushidocollective/brisk.git
   cd brisk
   ```

3. **Follow the specific README** for your chosen method

### Development Environment

For local development, use Docker Compose:

```bash
cd deployment/docker-compose
cp .env.example .env
# Edit .env with your configuration
docker-compose up -d
```

Access at: http://localhost

### Production Deployment

For production, choose Kubernetes or a managed platform:

**Kubernetes:**
```bash
helm install kasoku ./deployment/kubernetes/kasoku \
  --values deployment/kubernetes/kasoku/values-cloud.yaml \
  --set global.domain=kasoku.yourcompany.com
```

**Managed (Fly.io example):**
```bash
./deployment/one-click/fly/deploy.sh
```

---

## 🔒 Security Considerations

All deployment methods include:

- ✅ HTTPS/TLS encryption
- ✅ Secure secret management
- ✅ Non-root containers
- ✅ Network isolation
- ✅ Database authentication
- ✅ Health checks

**Production Checklist:**
- [ ] Use strong, randomly generated secrets
- [ ] Enable database backups
- [ ] Configure firewall rules
- [ ] Set up monitoring and alerts
- [ ] Use external secret management (Vault, AWS Secrets Manager)
- [ ] Enable audit logging
- [ ] Configure rate limiting
- [ ] Set up WAF (Web Application Firewall)

---

## 📊 Feature Comparison

| Feature | Docker Compose | Kubernetes | Railway | Render | Fly.io |
|---------|----------------|------------|---------|--------|--------|
| **Auto-scaling** | ❌ | ✅ | ✅ | ✅ | ✅ |
| **Multi-region** | ❌ | ✅ | ❌ | ❌ | ✅ |
| **Managed DB** | ❌ | Optional | ✅ | ✅ | ✅ |
| **Auto SSL** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Health Checks** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Monitoring** | Manual | ✅ | ✅ | ✅ | ✅ |
| **Backup** | Manual | Manual | Auto | Auto | Auto |
| **Cost** | Low | Medium | Low | Medium | Low |
| **Complexity** | Low | High | Low | Low | Medium |

---

## 📚 Additional Resources

### Documentation
- [Docker Compose Deployment](docker-compose/README.md)
- [Kubernetes Helm Chart](kubernetes/kasoku/README.md)
- [Railway Deployment](one-click/railway/README.md)
- [Render Deployment](one-click/render/README.md)
- [Fly.io Deployment](one-click/fly/README.md)

### Kasoku Documentation
- [Main README](../README.md)
- [CLI Documentation](../cmd/README.md)
- [Server API Documentation](../server/README.md)
- [Web Dashboard](../web/README.md)

### Support
- GitHub Issues: https://github.com/thebushidocollective/brisk/issues
- Discussions: https://github.com/thebushidocollective/brisk/discussions

---

## 🤝 Contributing

Found an issue with a deployment method? Have a suggestion?

1. Check existing issues
2. Create a new issue with:
   - Deployment method used
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details

Or submit a pull request with improvements!

---

## 📝 License

All deployment configurations are part of the Kasoku project and follow the same license.

See [LICENSE](../LICENSE) for details.
