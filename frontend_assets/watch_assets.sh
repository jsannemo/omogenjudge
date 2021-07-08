#!/usr/bin/env bash

assets_path=`realpath $(dirname $0)`

(cd $assets_path;
    npm install;
    rm -rf static;
    mkdir static;
    cp -r img static;
    rm -rf ../output/static;
    mkdir -p ../output/
    rm -rf ../output/frontend_assets/
    ln -s $assets_path/static ../output/frontend_assets;
    npm run watch;
)
