#!/bin/sh
set -e

# Run database seeding
/seed

# Start the server
exec /server
