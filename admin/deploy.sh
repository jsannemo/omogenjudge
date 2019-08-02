#!/usr/bin/env bash
set -e

admin/build.sh

sudo dpkg -i builds/omogenjudge-sandbox.deb
sudo dpkg -i builds/omogenjudge-local.deb
sudo dpkg -i builds/omogenjudge-master.deb
sudo dpkg -i builds/omogenjudge-frontend.deb
