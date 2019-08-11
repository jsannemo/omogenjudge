#!/usr/bin/env bash
set -e

scripts/build/packages.sh

sudo dpkg -i builds/omogenjudge-sandbox-dev.deb
sudo dpkg -i builds/omogenjudge-local-dev.deb
sudo dpkg -i builds/omogenjudge-master-dev.deb
sudo dpkg -i builds/omogenjudge-frontend-dev.deb
