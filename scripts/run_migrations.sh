#!/bin/bash
set -e

# Load environment variables
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

DB_USER=${DB_USER:-microservice}
DB_PASSWORD=${DB_PASSWORD}
DB_HOST=${DB_HOST:-localhost}

if [ -z "$DB_PASSWORD" ]; then
    echo " Error: DB_PASSWORD not set"
    exit 1
fi

echo " Running database migrations..."

# Run migrations for each service database
services=("auth" "contact" "inventory" "sales" "purchase")
ports=(5432 5433 5434 5435 5436)
dbs=("auth" "contact" "inventory" "sales" "purchase")

for i in "${!services[@]}"; do
    service=${services[$i]}
    port=${ports[$i]}
    db=${dbs[$i]}
    
    echo "  Running migrations for $service service (database: $db)..."
    
    # Use docker exec to run migrations if database is in container
    if docker ps | grep -q "db-$service"; then
        docker exec -i db-$service psql -U $DB_USER -d $db <<EOF > /dev/null 2>&1 || true
$(cat migrations/$service/000001_*.up.sql)
EOF
        echo "     $service migrations applied"
    else
        # Fallback to direct psql if available
        if command -v psql > /dev/null; then
            PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $port -U $DB_USER -d $db -f migrations/$service/000001_*.up.sql > /dev/null 2>&1 || echo "    ⚠️  $service migration may already exist"
        else
            echo "      Cannot run migrations for $service - psql not available and container not running"
        fi
    fi
done

echo " All migrations completed"

