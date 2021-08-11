
<p align="center"><img width="350" height="281" src="logo.png" alt="MBT logo"></p>

<br />

[![CircleCI](https://circleci.com/gh/SAP/cloud-mta-build-tool.svg?style=svg&circle-token=ecedd1dce3592adcd72ee4c61481972c32dcfad7)](https://circleci.com/gh/SAP/cloud-mta-build-tool)
[![Go Report Card](https://goreportcard.com/badge/github.com/SAP/cloud-mta-build-tool)](https://goreportcard.com/report/github.com/SAP/cloud-mta-build-tool)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3400/badge)](https://bestpractices.coreinfrastructure.org/projects/3400) 
[![Coverage Status](https://coveralls.io/repos/github/SAP/cloud-mta-build-tool/badge.svg?branch=cover)](https://coveralls.io/github/SAP/cloud-mta-build-tool?branch=cover)
![Beta](https://img.shields.io/badge/version-v1-green)
[![GitHub stars](https://img.shields.io/badge/contributions-welcome-orange.svg)](docs/docs/process.md)
[![dependentbot](https://api.dependabot.com/badges/status?host=github&repo=SAP/cloud-mta-build-tool)](https://dependabot.com/)
[![REUSE status](https://api.reuse.software/badge/github.com/SAP/cloud-mta-build-tool)](https://api.reuse.software/info/github.com/SAP/cloud-mta-build-tool)


## Description

#### Multi-Target Application

Before using this package, make sure that you are familiar with the multi-target application concept and terminology. For background and detailed information, see the [Multi-Target Application Model](https://www.sap.com/documents/2016/06/e2f618e4-757c-0010-82c7-eda71af511fa.html) guide. 

#### The Cloud MTA Build Tool Overview
The Cloud MTA Build Tool is a standalone command-line tool that builds a deployment-ready
multitarget application (MTA) archive `.mtar` file from the artifacts of an MTA project according to the project’s MTA
development descriptor (`mta.yaml` file) or from module build artifacts according to the MTA deployment descriptor (`mtad.yaml` file). Also, it provides commands for running intermediate build process steps; for example, the `mta.yaml` file validations, building a single module according to the configurations in the development descriptor, generating the deployment descriptor, and so on.


><b>For more information, see the [Cloud MTA Build Tool user guide](https://sap.github.io/cloud-mta-build-tool/)</b>

#### Demo

This demo shows the basic usage of the tool. For more advanced scenarios, follow the documentation.

<p align="center">
  <img src="./docs/demo.gif" width="100%">
</p>

#### The Cloud MTA Build Tool Images
The images can be used to build SAP Multitarget Applications (MTA) containing Java, Node.js, and Golang modules and **provided for CI env**. 
The images are hosted at [Github container registry](https://github.com/orgs/SAP/packages?tab=packages&q=mbtci-) and [Docker Hub registry](https://hub.docker.com/search?q=mbtci-&type=image).

**Note:** 
* For most cases, it's highly recommended to use the [alpine](https://hub.docker.com/r/devxci/mbtci-alpine) version, e.g. `docker pull devxci/mbtci-alpine` ,
this version is more light-weight and should be used in `production` env.
Using the `alpine` version gives the flexibility to add "per-scenario" the required set of tools. 
* The mbtci-alpine image is also hosted at [GitHub container](https://github.com/orgs/SAP/packages/container/package/mbtci-alpine).

##### How to pull the image

From the command line:
```
$ docker pull devxci/mbtci-alpine:latest
```
or
```
$ docker pull ghcr.io/sap/mbtci-alpine:latest
```
From Dockerfile as a base image:
```
FROM devxci/mbtci-alpine:latest
```
or
```
FROM ghcr.io/sap/mbtci-alpine:latest
```

##### How to use the image
On a Linux/Darwin machine you can run:
```
docker run -it --rm -v "$(pwd)/[proj-releative-path]:/project" devxci/mbtci-alpine:latest mbt build -p=cf -t [target-folder-name]
```
This will build an mtar file for SAP Cloud Platform (Cloud Foundry). The folder containing the project needs to be mounted into the image at /project.

<b>Note:</b> The parameter `-p=cf` can be omitted as the build for cloud foundry is the default build, this is an example of the MBT build parameters, for further commands see MBT docs.

##### How to build the image
First you need to copy the relevant Dockerfile according desired base image (alpine, sapjvm or sapmachine):
```
cp Dockerfile_<base image> Dockerfile
```
In case of base image sapjvm or sapmachine you need to replace NODE_VERSION_TEMPLATE with node version 12.18.3 or 14.17.0 in following line in the Dockerfile:
```
ARG NODE_VERSION=NODE_VERSION_TEMPLATE
```
Then you can build the image:
```
docker build -t devxci/mbtci .
```

##### The images provide:

- Cloud MTA Build Tool - 1.2.2
- Nodejs - 12.18.3 or 14.17.0
- Maven - 3.6.3
- Golang - 1.14.7
- Java - 8 or 11

The MTA Archive Builder delegates module builds to other native build tools. This image provides Node.js, Maven, Java, and Golang so the archive builder can delegate to these build technologies. In case other build tools are needed, <b>inherit</b> from this image and add more build tools.

##### License

Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
Please note that Docker images can contain other software which may be licensed under different licenses. This License file is also included in the Docker image. For any usage of built Docker images please make sure to check the licenses of the artifacts contained in the images.

## Contributions

Contributions are greatly appreciated.
If you want to contribute, follow [the guidelines](docs/docs/process.md).

## Support

Please follow our [issue template](https://github.com/SAP/cloud-mta-build-tool/issues/new/choose) on how to report an issue.
