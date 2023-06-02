#!/bin/sh

rm -rf node-js

cp -r ../cmd/testdata/mta/node-js ./

cd node-js 

npm install --production

cd -

node test-micromatch.js

go run test_micromatch_wrapper_js.go

go run test_micromatch_wrapper_bin.go