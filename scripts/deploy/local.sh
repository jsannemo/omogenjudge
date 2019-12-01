#!/usr/bin/env bash
set -e

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd $SCRIPT_DIR/../../

scripts/build/packages.sh

sudo dpkg -i builds/omogenjudge-local-dev.deb
