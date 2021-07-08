#!/usr/bin/env bash
set -e

sudo service omogenjudge-queue stop || true

sudo -u postgres psql omogenjudge -c "DELETE FROM submission;"

sudo service omogenjudge-queue start || true
