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

bazel build admin/addproblem:addproblem && bazel-bin/admin/addproblem/linux_amd64_stripped/addproblem docs/examples/problems/addition
psql omogenjudge -c "INSERT INTO account(username, password_hash, full_name, email) VALUES('test', '\$2a\$10\$r8xXriU.jnygztki.9eCv.C91FgU4BXnK/4Kl087v8RWsfGW0wcwW', 'Test Test', 'test@test-email.invalid');"
