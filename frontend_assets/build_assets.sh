#!/usr/bin/env bash

assets_path=`dirname $0`

(cd $assets_path;
    npm install;
    rm -rf static;
    npm run build;
    cp -r img static;
    rm -rf ../output/frontend_assets;
    mv static ../output/frontend_assets;
)
