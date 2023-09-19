#!/bin/sh

set -e # the script will exit immediately if command return no zero status

echo "start the app"
exec "$@" # take all parameters passed to the script and run it
