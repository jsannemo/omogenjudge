#!/usr/bin/env bash

set -e

base_path=`dirname -- "$0"`/..
judgehost_path=$base_path/judgehost
packaging_path=$base_path/packaging

(cd $judgehost_path && bazel build ...)
rm -f $packaging_path/omogenjudge-queue.deb
cp $judgehost_path/bazel-bin/queue/deb/omogenjudge-queue.deb $packaging_path
chmod 664 $packaging_path/omogenjudge-queue.deb
