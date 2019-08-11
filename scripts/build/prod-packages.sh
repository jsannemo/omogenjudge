#!/usr/bin/env bash
set -e

bazel build -c opt {frontend,localjudge,masterjudge,sandbox}/deb/...
mkdir -p builds

cp -f bazel-bin/frontend/deb/omogenjudge-frontend.deb builds/omogenjudge-frontend.deb
cp -f bazel-bin/localjudge/deb/omogenjudge-local.deb builds/omogenjudge-local.deb
cp -f bazel-bin/masterjudge/deb/omogenjudge-master.deb builds/omogenjudge-master.deb
cp -f bazel-bin/sandbox/deb/omogenjudge-sandbox.deb builds/omogenjudge-sandbox.deb
