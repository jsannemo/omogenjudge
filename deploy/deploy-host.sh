#!/usr/bin/env bash

set -e

base_path=`dirname -- "$0"`/..
(cd $base_path;
./packaging/build-host.sh;
cp ./packaging/omogenjudge-host.deb ./deploy/files/packages/omogenjudge-host.deb
)

(cd $base_path; poetry run ansible-playbook -i deploy/hosts deploy/host.yml)
