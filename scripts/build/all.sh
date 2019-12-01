#!/usr/bin/env bash

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd $SCRIPT_DIR/../../

bazel build {admin,filehandler,frontend,localjudge,masterjudge,problemtools,runner,sandbox,schema,storage,util,vendor}/...
