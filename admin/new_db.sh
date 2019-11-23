#!/usr/bin/env bash
set -e

dropdb omogenjudge --if-exists
dropuser omogenjudge --if-exists
psql -c "CREATE USER omogenjudge WITH PASSWORD 'omogenjudge';"
createdb omogenjudge
psql omogenjudge -f schema/schema.sql

