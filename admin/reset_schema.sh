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

./admin/new_db.sh

sudo service omogenjudge-local start
sudo service omogenjudge-master start
if [[ $nofrontend -ne 1 ]]; then
  sudo service omogenjudge-frontend start
fi
bazel build admin/addproblem:addproblem && bazel-bin/admin/addproblem/linux_amd64_stripped/addproblem docs/example-problems/addition
bazel build admin/addcourse:addcourse && bazel-bin/admin/addcourse/linux_amd64_stripped/addcourse docs/example-courses/python-intro
