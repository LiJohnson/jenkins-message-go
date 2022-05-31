#!/bin/bash 

mkdir dist || echo "dist exist"
rm -rf dist/*

# npm i -g html-minifier-terser
cat home.html \
 | html-minifier-terser --collapse-whitespace --remove-comments --minify-css true --minify-js true \
   -o dist/home.html

build_time=$(date +%Y%m%dT%H%M%S)
git_hash=$(git rev-parse HEAD)

go build -ldflags "-X main.buildTime=${build_time}  -X main.gitHash=${git_hash}"  -o dist/