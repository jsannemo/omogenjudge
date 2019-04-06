#!/bin/bash

RELEASE_VERSION=0.3.3
curl -L https://github.com/grailbio/bazel-compilation-database/archive/${RELEASE_VERSION}.tar.gz | tar -xz
bazel-compilation-database-${RELEASE_VERSION}/generate.sh
