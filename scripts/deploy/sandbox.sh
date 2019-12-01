#!/usr/bin/env bash
set -e

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd $SCRIPT_DIR/../../

scripts/build/packages.sh

sudo dpkg -i builds/omogenjudge-sandbox-dev.deb

sudo service omogenjudge-local restart
sudo service omogenjudge-master restart
