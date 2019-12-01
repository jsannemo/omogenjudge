#!/usr/bin/env bash
set -e

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd $SCRIPT_DIR/../../

scripts/build/packages.sh

sudo dpkg -i builds/omogenjudge-sandbox-dev.deb
sudo dpkg -i builds/omogenjudge-local-dev.deb
sudo dpkg -i builds/omogenjudge-master-dev.deb
sudo dpkg -i builds/omogenjudge-frontend-dev.deb
