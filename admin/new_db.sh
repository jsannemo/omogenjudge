#!/usr/bin/env bash
set -e

sudo service omogenjudge-queue stop
sudo service omogenjudge-host stop

dropdb omogenjudge --if-exists
dropuser omogenjudge --if-exists
dropuser omogenhost --if-exists
psql -c "CREATE USER omogenjudge WITH PASSWORD 'omogenjudge';"
psql -c "CREATE USER omogenhost WITH PASSWORD 'omogenhost';"
createdb omogenjudge

sudo service omogenjudge-host start
sudo service omogenjudge-queue start
