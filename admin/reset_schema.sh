#!/usr/bin/env bash
set -e

dropdb omogenjudge --if-exists
dropuser omogenjudge --if-exists
psql -c "CREATE USER omogenjudge WITH PASSWORD 'omogenjudge';"
createdb omogenjudge
psql omogenjudge -f schema/schema.sql

psql omogenjudge -c "INSERT INTO account(username, password_hash) VALUES('test', '\$2a\$10\$r8xXriU.jnygztki.9eCv.C91FgU4BXnK/4Kl087v8RWsfGW0wcwW');"
bazel build admin/addproblem:addproblem && bazel-bin/admin/addproblem/linux_amd64_stripped/addproblem docs/example-problems/addition
