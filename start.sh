#!/bin/sh

set -e

# echo "Current working directory: $(pwd)"
# echo "Listing contents of current directory:"
# ls -la

# echo "Listing contents of db/migration:"
# ls -la /app/db/migration

echo "Running db migration..."
source /app/app.env
/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

exec "$@"