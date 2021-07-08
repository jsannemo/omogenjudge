#!/usr/bin/env bash
set -e

dropdb omogenjudge --if-exists
dropuser omogenjudge --if-exists
dropuser omogenhost --if-exists
psql -c "CREATE USER omogenjudge WITH PASSWORD 'omogenjudge';"
psql -c "CREATE USER omogenhost WITH PASSWORD 'omogenhost';"
createdb omogenjudge
