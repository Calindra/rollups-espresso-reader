#!/bin/sh

# Wait until PostgreSQL is ready
# wait_for_postgres() {
#   echo "Waiting for PostgreSQL to become available..."
#   until PGPASSWORD="$POSTGRES_PASSWORD" psql -h "$DB_SERVICE_NAME" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q' 2>/dev/null; do
#     echo "PostgreSQL not ready yet - waiting..."
#     sleep 2
#   done
#   echo "PostgreSQL is ready!"
# }

# Configure database connection
export CARTESI_DATABASE_CONNECTION="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$DB_SERVICE_NAME:5432/$POSTGRES_DB?sslmode=disable"

# Wait for PostgreSQL to be available
# wait_for_postgres

# Run rollups-node migrations
echo "Running rollups-node migrations..."
cd /usr/src/app/rollups-node
eval $(make env)
export CARTESI_DATABASE_CONNECTION="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$DB_SERVICE_NAME:5432/$POSTGRES_DB?sslmode=disable"
make migrate

# Return to root directory and run espresso-reader migrations
echo "Running espresso-reader migrations..."
cd /usr/src/app
eval $(make env)
export CARTESI_DATABASE_CONNECTION="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$DB_SERVICE_NAME:5432/$POSTGRES_DB?sslmode=disable"
make migrate

# Run migrations and generate DB code
if [ "${SKIP_GEN_MIGRATIONS:-false}" = "false" ]; then
    make generate-db
fi

# Keep container running if required
if [ "${KEEP_RUNNING:-false}" != "false" ]; then
    echo "Migrations completed. Container running..."
    tail -f /dev/null
else
    echo "Migrations completed."
fi