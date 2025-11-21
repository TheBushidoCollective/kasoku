#!/bin/bash
set -e

# Kasoku Database Backup Script
# This script backs up the PostgreSQL database and storage volumes

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
BACKUP_DIR="${BACKUP_DIR:-./backups}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.production.yml}"

echo -e "${GREEN}🔄 Kasoku Backup - $(date)${NC}"
echo ""

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo -e "${RED}Error: .env file not found${NC}"
    exit 1
fi

# Backup PostgreSQL database
echo -e "${YELLOW}📦 Backing up PostgreSQL database...${NC}"
docker-compose -f "$COMPOSE_FILE" exec -T postgres pg_dump -U kasoku kasoku | gzip > "$BACKUP_DIR/postgres_${TIMESTAMP}.sql.gz"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Database backup complete: postgres_${TIMESTAMP}.sql.gz${NC}"
else
    echo -e "${RED}✗ Database backup failed${NC}"
    exit 1
fi

# Backup storage volume (cache artifacts)
echo -e "${YELLOW}📦 Backing up storage volume...${NC}"
docker run --rm \
    -v "$(pwd)_storage_data:/data" \
    -v "$BACKUP_DIR:/backup" \
    alpine tar czf "/backup/storage_${TIMESTAMP}.tar.gz" -C /data .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Storage backup complete: storage_${TIMESTAMP}.tar.gz${NC}"
else
    echo -e "${RED}✗ Storage backup failed${NC}"
    exit 1
fi

# Backup .env file (important for restoration)
echo -e "${YELLOW}📦 Backing up configuration...${NC}"
cp .env "$BACKUP_DIR/env_${TIMESTAMP}.backup"
echo -e "${GREEN}✓ Configuration backup complete${NC}"

# Calculate backup sizes
DB_SIZE=$(du -h "$BACKUP_DIR/postgres_${TIMESTAMP}.sql.gz" | cut -f1)
STORAGE_SIZE=$(du -h "$BACKUP_DIR/storage_${TIMESTAMP}.tar.gz" | cut -f1)

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✓ Backup completed successfully${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Database:      $DB_SIZE (postgres_${TIMESTAMP}.sql.gz)"
echo "Storage:       $STORAGE_SIZE (storage_${TIMESTAMP}.tar.gz)"
echo "Configuration: env_${TIMESTAMP}.backup"
echo "Location:      $BACKUP_DIR"
echo ""

# Cleanup old backups
if [ "$RETENTION_DAYS" -gt 0 ]; then
    echo -e "${YELLOW}🧹 Cleaning up backups older than $RETENTION_DAYS days...${NC}"

    # Find and delete old backups
    OLD_BACKUPS=$(find "$BACKUP_DIR" -name "postgres_*.sql.gz" -mtime +$RETENTION_DAYS)
    OLD_STORAGE=$(find "$BACKUP_DIR" -name "storage_*.tar.gz" -mtime +$RETENTION_DAYS)
    OLD_ENV=$(find "$BACKUP_DIR" -name "env_*.backup" -mtime +$RETENTION_DAYS)

    if [ -n "$OLD_BACKUPS" ] || [ -n "$OLD_STORAGE" ] || [ -n "$OLD_ENV" ]; then
        find "$BACKUP_DIR" -name "postgres_*.sql.gz" -mtime +$RETENTION_DAYS -delete
        find "$BACKUP_DIR" -name "storage_*.tar.gz" -mtime +$RETENTION_DAYS -delete
        find "$BACKUP_DIR" -name "env_*.backup" -mtime +$RETENTION_DAYS -delete
        echo -e "${GREEN}✓ Old backups removed${NC}"
    else
        echo "No old backups to remove"
    fi
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "To restore from this backup:"
echo ""
echo "  # Restore database:"
echo "  gunzip < $BACKUP_DIR/postgres_${TIMESTAMP}.sql.gz | \\"
echo "    docker-compose -f $COMPOSE_FILE exec -T postgres psql -U kasoku kasoku"
echo ""
echo "  # Restore storage:"
echo "  docker run --rm -v $(pwd)_storage_data:/data -v $BACKUP_DIR:/backup \\"
echo "    alpine tar xzf /backup/storage_${TIMESTAMP}.tar.gz -C /data"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
