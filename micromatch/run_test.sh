#!/bin/sh

rm -rf ./micromatch-wrapper-linux ./micromatch-wrapper-macos ./micromatch-wrapper-win.exe

npm install

pkg ./

go run test_micromatch_wrapper_js.go

go run test_micromatch_wrapper_bin.go

cp ./micromatch-wrapper-win.exe $GOPATH/bin/

micromatch-wrapper-win.exe -h