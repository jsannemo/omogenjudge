#!/usr/bin/env bash
set -e

nofrontend=0

while [[ "$1" =~ ^- && ! "$1" == "--" ]]; do case $1 in
  --nofrontend )
    nofrontend=1
esac; shift; done


sudo service omogenjudge-local stop
sudo service omogenjudge-master stop
sudo service omogenjudge-frontend stop
dropdb omogenjudge --if-exists
dropuser omogenjudge --if-exists
psql -c "CREATE USER omogenjudge WITH PASSWORD 'omogenjudge';"
createdb omogenjudge
psql omogenjudge -f schema/schema.sql

psql omogenjudge -c "INSERT INTO account(username, password_hash) VALUES('test', '\$2a\$10\$r8xXriU.jnygztki.9eCv.C91FgU4BXnK/4Kl087v8RWsfGW0wcwW');"
sudo service omogenjudge-local start
sudo service omogenjudge-master start
if [[ $nofrontend -ne 1 ]]; then
  sudo service omogenjudge-frontend start
fi
bazel build admin/addproblem:addproblem && bazel-bin/admin/addproblem/linux_amd64_stripped/addproblem docs/example-problems/addition
bazel build admin/addcourse:addcourse && bazel-bin/admin/addcourse/linux_amd64_stripped/addcourse docs/example-courses/python-intro
