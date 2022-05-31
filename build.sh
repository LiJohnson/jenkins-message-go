#!/bin/bash 
# by lcs
# 2022-05-31

mkdir dist > /dev/null 2>&1
rm -rf dist/*

# npm i -g html-minifier-terser
cat home.html \
 | html-minifier-terser --collapse-whitespace --remove-comments --minify-css true --minify-js true \
   -o dist/home.html

build_time=$(date +%Y%m%dT%H%M%S)
git_hash=$(git rev-parse HEAD)

go build -ldflags "-X main.buildTime=${build_time}  -X main.gitHash=${git_hash}"  -o dist/

cd dist

zip dist.zip ./*