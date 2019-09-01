#!/usr/bin/env bash
set -e

scripts/build/packages.sh

sudo dpkg -i builds/omogenjudge-frontend-dev.deb
