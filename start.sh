#!/bin/sh

set -e # the script will exit immediately if command return no zero status

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@" # take all parameters passed to the script and run it
