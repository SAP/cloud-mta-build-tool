#!/bin/sh

# clean env
rm -rf ./micromatch-wrapper-linux ./micromatch-wrapper-macos ./micromatch-wrapper-win.exe
rm -rf node_modules

# install pkg
npm install -g pkg
pkg --version

# build micromatch wrapper
npm install
pkg ./

# install and test micromatch wrapper
cp ./micromatch-wrapper-win.exe $GOPATH/bin/micromatch-wrapper.exe
micromatch-wrapper.exe -h

# clean and copy micromatch wrapper to target path for release to GitHub
rm -rf ./Linux/* ./Darwin/* ./Windows/*
mv ./micromatch-wrapper-linux ./Linux/micromatch-wrapper
mv ./micromatch-wrapper-macos ./Darwin/micromatch-wrapper
mv ./micromatch-wrapper-win.exe ./Windows/micromatch-wrapper.exe