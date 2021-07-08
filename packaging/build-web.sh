#!/usr/bin/env bash

set -e

base_path=`dirname -- "$0"`/..
(cd $base_path;
rm -rf output;
mkdir output;
./frontend_assets/build_assets.sh;
poetry install
poetry run python manage.py collectstatic -c --no-input;
cp -r packaging/web/ output/webdeb/;
cp -r omogenjudge output/webdeb/;
cp -r bin output/webdeb/;
cp pyproject.toml output/webdeb/;
cp -r output/static output/webdeb/;
)

(cd $base_path/output/webdeb; dpkg-buildpackage -us -uc -b)

mv $base_path/output/omogenjudge-web*.deb $base_path/packaging/omogenjudge-web.deb

