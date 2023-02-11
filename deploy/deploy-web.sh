#!/usr/bin/env bash

set -e

base_path=`dirname -- "$0"`/..
(cd $base_path;
./packaging/build-web.sh;
mkdir -p ./deploy/files/packages;
cp ./packaging/omogenjudge-web.deb ./deploy/files/packages/omogenjudge-web.deb
)

(cd $base_path; poetry run ansible-playbook -i deploy/hosts deploy/web.yml)
