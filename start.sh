#!/bin/sh

set -e

echo "run db migrations"

source /app/app.env

/app/migrate -path /app/migrations -database "$DB_DSN" -verbose up

echo "start app"

exec "$@"