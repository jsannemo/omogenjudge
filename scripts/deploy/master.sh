#!/usr/bin/env bash
set -e

scripts/build/packages.sh

sudo dpkg -i builds/omogenjudge-master-dev.deb
