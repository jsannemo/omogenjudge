#!/usr/bin/env bash

set -e

base_path=`dirname -- "$0"`/..
(cd $base_path;
./packaging/build-queue.sh;
cp ./packaging/omogenjudge-queue.deb ./deploy/files/packages/omogenjudge-queue.deb
)

(cd $base_path; poetry run ansible-playbook -i deploy/hosts deploy/queue.yml)
