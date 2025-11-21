#!/bin/bash
set -e

echo "🚀 Deploying Kasoku to Fly.io"
echo "=============================="
echo ""

# Check if flyctl is installed
if ! command -v fly &> /dev/null; then
    echo "❌ Fly CLI is not installed"
    echo "Install it with: curl -L https://fly.io/install.sh | sh"
    exit 1
fi

# Check if logged in
if ! fly auth whoami &> /dev/null; then
    echo "❌ Not logged in to Fly.io"
    echo "Run: fly auth login"
    exit 1
fi

echo "✅ Fly CLI detected"
echo ""

# Prompt for configuration
read -p "Enter your preferred region (e.g., sjc, iad, lhr): " REGION
REGION=${REGION:-sjc}

read -p "Create PostgreSQL database? (y/n): " CREATE_DB
read -p "Enable Stripe billing? (y/n): " ENABLE_BILLING

echo ""
echo "Configuration:"
echo "  Region: $REGION"
echo "  Create DB: $CREATE_DB"
echo "  Billing: $ENABLE_BILLING"
echo ""

read -p "Proceed with deployment? (y/n): " CONFIRM
if [[ ! $CONFIRM =~ ^[Yy]$ ]]; then
    echo "Deployment cancelled"
    exit 0
fi

# Step 1: Create and deploy PostgreSQL (if requested)
if [[ $CREATE_DB =~ ^[Yy]$ ]]; then
    echo ""
    echo "📦 Step 1: Creating PostgreSQL database..."
    fly postgres create \
        --name kasoku-db \
        --region "$REGION" \
        --initial-cluster-size 1 \
        --vm-size shared-cpu-1x \
        --volume-size 10

    echo ""
    echo "⚠️  Save the database connection string above!"
    read -p "Press Enter to continue..."
fi

# Get database connection string
if [[ $CREATE_DB =~ ^[Yy]$ ]]; then
    read -p "Enter the DATABASE_URL from above: " DB_DSN
else
    read -p "Enter your external DATABASE_URL: " DB_DSN
fi

# Step 2: Deploy server
echo ""
echo "🔧 Step 2: Deploying Kasoku Server..."

# Create app if it doesn't exist
if ! fly apps list | grep -q "kasoku-server"; then
    fly apps create kasoku-server
fi

# Create storage volume
if ! fly volumes list --app kasoku-server | grep -q "kasoku_storage"; then
    echo "Creating storage volume..."
    fly volumes create kasoku_storage \
        --region "$REGION" \
        --size 10 \
        --app kasoku-server
fi

# Generate JWT secret
JWT_SECRET=$(openssl rand -hex 32)

# Set server secrets
echo "Setting server secrets..."
fly secrets set \
    JWT_SECRET="$JWT_SECRET" \
    DB_DSN="$DB_DSN" \
    --app kasoku-server

# Optional: Stripe billing
if [[ $ENABLE_BILLING =~ ^[Yy]$ ]]; then
    echo ""
    read -p "Enter Stripe Secret Key: " STRIPE_SECRET_KEY
    read -p "Enter Stripe Webhook Secret: " STRIPE_WEBHOOK_SECRET
    read -p "Enter Individual Price ID: " STRIPE_INDIVIDUAL_PRICE_ID
    read -p "Enter Team Price ID: " STRIPE_TEAM_PRICE_ID

    fly secrets set \
        STRIPE_SECRET_KEY="$STRIPE_SECRET_KEY" \
        STRIPE_WEBHOOK_SECRET="$STRIPE_WEBHOOK_SECRET" \
        STRIPE_INDIVIDUAL_PRICE_ID="$STRIPE_INDIVIDUAL_PRICE_ID" \
        STRIPE_TEAM_PRICE_ID="$STRIPE_TEAM_PRICE_ID" \
        --app kasoku-server
fi

# Deploy server
echo "Deploying server..."
fly deploy \
    --config deployment/one-click/fly/fly.toml \
    --app kasoku-server

# Get server URL
SERVER_URL=$(fly apps list | grep kasoku-server | awk '{print "https://" $1 ".fly.dev"}')
echo "✅ Server deployed: $SERVER_URL"

# Step 3: Deploy web
echo ""
echo "🌐 Step 3: Deploying Kasoku Web..."

# Create app if it doesn't exist
if ! fly apps list | grep -q "kasoku-web"; then
    fly apps create kasoku-web
fi

# Set web environment variables
echo "Setting web environment variables..."
fly secrets set \
    NEXT_PUBLIC_API_URL="$SERVER_URL" \
    --app kasoku-web

# Deploy web
echo "Deploying web..."
fly deploy \
    --config deployment/one-click/fly/fly-web.toml \
    --app kasoku-web

# Get web URL
WEB_URL=$(fly apps list | grep kasoku-web | awk '{print "https://" $1 ".fly.dev"}')
echo "✅ Web deployed: $WEB_URL"

# Summary
echo ""
echo "🎉 Deployment Complete!"
echo "======================"
echo ""
echo "Your Kasoku installation is ready:"
echo ""
echo "  Web Dashboard: $WEB_URL"
echo "  API Server:    $SERVER_URL"
echo ""
echo "Next steps:"
echo "  1. Visit $WEB_URL to access Kasoku"
echo "  2. Configure custom domain (optional):"
echo "     fly certs add your-domain.com --app kasoku-web"
echo "  3. Monitor logs:"
echo "     fly logs --app kasoku-server"
echo "     fly logs --app kasoku-web"
echo ""
echo "For help, visit: https://github.com/thebushidocollective/brisk"
echo ""
