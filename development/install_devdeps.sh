#!/usr/bin/env bash
set -e

echo "Installing Python for boostrapping"
sudo apt install python3 python3-dev

echo "Installing poetry"
curl -sSL https://install.python-poetry.org | python3 -

echo "Installing Node"
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

echo "Installing postgres"
sudo apt install postgresql libpq-dev

echo "Installing Python build dependencies"
sudo apt install build-essential

echo "Installing problemtools dependencies"
sudo apt install automake libgmp-dev libboost-regex-dev

echo "Installing bazel"
sudo npm install -g @bazel/bazelisk

echo "Installing Python dependencies"
poetry install
