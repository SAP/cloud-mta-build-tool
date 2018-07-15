#!/bin/bash
# This script build the mbt for various env
# Create build with OS artifact's which need to put under the bin file as executable bin

basedir=$(cd -- "${BASH_SOURCE%/*}" && pwd) || exit

rm -rf -- "$basedir/"build   || exit
# Create build folder to drop the artifacts
mkdir -p -- "$basedir/"build  || exit

build() (GOOS=$1 GOARCH=$2 exec go build -o "$basedir/build/$3")

## build for specified OS
build darwin  amd64 mit             || exit
build linux   amd64 mit_linux       || exit
build windows amd64 mit_win64       || exit

# copy the new artifacts to the bin folder to save time
cp "$basedir/"build/mit $GOPATH/bin/




