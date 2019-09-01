#!/usr/bin/env bash
set -e

scripts/build/packages.sh

sudo dpkg -i builds/omogenjudge-sandbox-dev.deb

sudo service omogenjudge-local restart
sudo service omogenjudge-master restart
