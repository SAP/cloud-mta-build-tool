#!/bin/sh

git tag --delete release
git push https://github.com/SAP/cloud-mta-build-tool.git --delete release

git tag release
git push -u origin release

