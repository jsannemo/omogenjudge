#!/usr/bin/env bash
set -e

bazel build -c opt {frontend,localjudge,masterjudge,sandbox}/deb/...
mkdir builds
cp bazel-bin/{frontend,localjudge,masterjudge,sandbox}/deb/*.deb builds
