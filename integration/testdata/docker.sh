#!/usr/bin/env sh

export NODE_TEST_FOLDER=node
export MAVEN_TEST_FOLDER=maven

docker build -t mbt "$(pwd)/.." || echo "error occured while building MTA project"
docker run -it --rm -v "$(pwd)/${NODE_TEST_FOLDER}:/project" cmbt3 mbt build -p=xsa
