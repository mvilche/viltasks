#!/bin/bash -e
#export GOARCH=amd64
#export GOOS=linux
rm -rf build/
revel build .  build/
cp -rf conf/ build/
mkdir -p build/database