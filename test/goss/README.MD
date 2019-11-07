# Testing the Container with DGOSS

[GOSS](https://github.com/aelsabbahy/goss) is a *server validation tool* and together with [DGOSS](https://github.com/aelsabbahy/goss/tree/master/extras/dgoss) it can be used to validate Docker conteiners.

## Execute the Tests

1. Build the image in the project's root directory:
```shell
docker build -t devxci/mbtci .
```
2. 
```shell
GOSS_FILES_PATH=./test/goss dgoss run --interactive --volume $PWD:/project devxci/mbtci
```

## Edit the Tests

To edit the tests start the edit mode:
```shell
GOSS_FILES_PATH=./test/goss dgoss edit --interactive --volume $PWD:/project devxci/mbtci
```
Add further test cases, see `goss add --help` for further informations.

## Execute the Tests with Docker-Compose

The tests can also be executed with *docker-compose* which makes it obsolete to install GOSS/DGOSS.

Execute the *docker-compose* file:
```
docker-compose --file docker-compose.test.yml up --build
```

The `--build` option forces the *image under test* to be rebuilt.

## Docker Hub Autotest

The Docker build on *Docker Hub* can automatically test the image that is built. This happend via a *docker-compose* file named `docker-compose.test.yml` that contains a *service* `sut`.
Further informations can be found in [the official docs](https://docs.docker.com/docker-hub/builds/automated-testing/).

