#!/bin/sh

# build and install micromatch wrapper
rm -rf ./micromatch-wrapper-linux ./micromatch-wrapper-macos ./micromatch-wrapper-win.exe

npm install

pkg ./

cp ./micromatch-wrapper-win.exe $GOPATH/bin/micromatch-wrapper.exe

micromatch-wrapper.exe -h