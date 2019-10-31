#!/usr/bin/env bash

export NODE_TEST_FOLDER=node
export MAVEN_TEST_FOLDER=maven

docker build -t mbtci "$(pwd)/.." || echo "error occured while building MTA project"
docker run -it --rm -v "$(pwd)/${NODE_TEST_FOLDER}:/project" mbtci mbt build -p=xsa

docker run -it --rm -v "$(pwd)/${MAVEN_TEST_FOLDER}:/project" mbtci mbt build -p=cf