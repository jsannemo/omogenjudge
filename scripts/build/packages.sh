#!/usr/bin/env bash
set -e

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd $SCRIPT_DIR/../../

bazel build {frontend,localjudge,masterjudge,sandbox}/deb/...
mkdir -p builds

cp -f bazel-bin/frontend/deb/omogenjudge-frontend.deb builds/omogenjudge-frontend-dev.deb
cp -f bazel-bin/localjudge/deb/omogenjudge-local.deb builds/omogenjudge-local-dev.deb
cp -f bazel-bin/masterjudge/deb/omogenjudge-master.deb builds/omogenjudge-master-dev.deb
cp -f bazel-bin/sandbox/deb/omogenjudge-sandbox.deb builds/omogenjudge-sandbox-dev.deb
