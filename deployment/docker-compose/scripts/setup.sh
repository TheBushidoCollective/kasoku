#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "╔════════════════════════════════════════════════════╗"
echo "║                                                    ║"
echo "║     Kasoku Production Setup                        ║"
echo "║     Self-hosted cache acceleration                ║"
echo "║                                                    ║"
echo "╚════════════════════════════════════════════════════╝"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run as root (or use sudo)${NC}"
    exit 1
fi

# Check Docker installation
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}📦 Docker not found. Installing Docker...${NC}"
    curl -fsSL https://get.docker.com | sh
    echo -e "${GREEN}✓ Docker installed${NC}"
else
    echo -e "${GREEN}✓ Docker is already installed${NC}"
fi

# Check Docker Compose installation
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo -e "${YELLOW}📦 Docker Compose not found. Installing...${NC}"
    # Docker Compose V2 is usually included with Docker now
    echo -e "${RED}Please install Docker Compose manually:${NC}"
    echo "https://docs.docker.com/compose/install/"
    exit 1
else
    echo -e "${GREEN}✓ Docker Compose is installed${NC}"
fi

# Navigate to deployment directory
cd "$(dirname "$0")/.."

# Check if .env already exists
if [ -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file already exists${NC}"
    read -p "Do you want to overwrite it? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Using existing .env file"
        exit 0
    fi
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Configuration Setup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Prompt for domain
read -p "Enter your domain (e.g., kasoku.example.com): " DOMAIN
if [ -z "$DOMAIN" ]; then
    echo -e "${RED}Domain is required!${NC}"
    exit 1
fi

# Generate secure secrets
echo -e "${GREEN}🔐 Generating secure secrets...${NC}"
JWT_SECRET=$(openssl rand -hex 32)
DB_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-32)

# Ask about storage type
echo ""
echo "Storage options:"
echo "  1) Local (recommended for self-hosted)"
echo "  2) S3 (for cloud/production at scale)"
read -p "Choose storage type (1 or 2): " STORAGE_CHOICE

if [ "$STORAGE_CHOICE" = "2" ]; then
    STORAGE_TYPE="s3"
    read -p "Enter S3 bucket name: " S3_BUCKET
    read -p "Enter S3 region (e.g., us-east-1): " S3_REGION
    read -p "Enter AWS Access Key ID: " AWS_ACCESS_KEY_ID
    read -p "Enter AWS Secret Access Key: " AWS_SECRET_ACCESS_KEY
else
    STORAGE_TYPE="local"
    S3_BUCKET=""
    S3_REGION=""
    AWS_ACCESS_KEY_ID=""
    AWS_SECRET_ACCESS_KEY=""
fi

# Ask about billing
echo ""
read -p "Enable Stripe billing? (y/N): " -n 1 -r ENABLE_BILLING
echo
if [[ $ENABLE_BILLING =~ ^[Yy]$ ]]; then
    read -p "Enter Stripe Secret Key: " STRIPE_SECRET_KEY
    read -p "Enter Stripe Webhook Secret: " STRIPE_WEBHOOK_SECRET
    read -p "Enter Stripe Individual Price ID: " STRIPE_INDIVIDUAL_PRICE_ID
    read -p "Enter Stripe Team Price ID: " STRIPE_TEAM_PRICE_ID
    STRIPE_SUCCESS_URL="https://$DOMAIN/dashboard?success=true"
    STRIPE_CANCEL_URL="https://$DOMAIN/pricing"
    STRIPE_RETURN_URL="https://$DOMAIN/dashboard"
else
    STRIPE_SECRET_KEY=""
    STRIPE_WEBHOOK_SECRET=""
    STRIPE_INDIVIDUAL_PRICE_ID=""
    STRIPE_TEAM_PRICE_ID=""
    STRIPE_SUCCESS_URL=""
    STRIPE_CANCEL_URL=""
    STRIPE_RETURN_URL=""
fi

# Create .env file
cat > .env <<EOF
# Kasoku Production Configuration
# Generated on $(date)

# Domain
DOMAIN=$DOMAIN
HTTP_PORT=80
HTTPS_PORT=443

# Database
DB_PASSWORD=$DB_PASSWORD

# Security
JWT_SECRET=$JWT_SECRET

# Storage
STORAGE_TYPE=$STORAGE_TYPE
S3_BUCKET=$S3_BUCKET
S3_REGION=$S3_REGION
AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY

# Stripe Billing
STRIPE_SECRET_KEY=$STRIPE_SECRET_KEY
STRIPE_WEBHOOK_SECRET=$STRIPE_WEBHOOK_SECRET
STRIPE_INDIVIDUAL_PRICE_ID=$STRIPE_INDIVIDUAL_PRICE_ID
STRIPE_TEAM_PRICE_ID=$STRIPE_TEAM_PRICE_ID
STRIPE_SUCCESS_URL=$STRIPE_SUCCESS_URL
STRIPE_CANCEL_URL=$STRIPE_CANCEL_URL
STRIPE_RETURN_URL=$STRIPE_RETURN_URL
EOF

# Set proper permissions
chmod 600 .env

echo ""
echo -e "${GREEN}✓ Configuration saved to .env${NC}"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Next Steps"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Point your domain ($DOMAIN) to this server's IP address"
echo "   Add an A record: $DOMAIN → $(curl -s ifconfig.me)"
echo ""
echo "2. Wait for DNS to propagate (5-30 minutes)"
echo "   Check with: dig $DOMAIN"
echo ""
echo "3. Start Kasoku:"
echo "   docker-compose -f docker-compose.production.yml up -d"
echo ""
echo "4. Check logs:"
echo "   docker-compose -f docker-compose.production.yml logs -f"
echo ""
echo "5. Visit your installation:"
echo "   https://$DOMAIN"
echo ""
echo -e "${YELLOW}⚠️  SSL certificates will be automatically provisioned on first access${NC}"
echo -e "${YELLOW}⚠️  This may take 1-2 minutes${NC}"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
