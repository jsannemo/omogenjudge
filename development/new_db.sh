#!/usr/bin/env bash
set -e

sudo service omogenjudge-queue stop || true
sudo service omogenjudge-host stop || true

sudo -u postgres dropdb omogenjudge --if-exists
sudo -u postgres dropuser omogenjudge --if-exists
sudo -u postgres dropuser omogenhost --if-exists
sudo -u postgres psql -c "CREATE USER omogenjudge WITH PASSWORD 'omogenjudge';"
sudo -u postgres createdb omogenjudge

sudo service omogenjudge-host start || true
sudo service omogenjudge-queue start || true

poetry run python manage.py migrate
