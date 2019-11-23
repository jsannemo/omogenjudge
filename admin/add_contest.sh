#!/usr/bin/env bash

bazel build admin/addcontest && bazel-bin/admin/addcontest/linux_amd64_stripped/addcontest $@
