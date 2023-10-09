
<p align="center"><img width="350" height="281" src="logo.png" alt="MBT logo"></p>

<br />

[![CircleCI](https://circleci.com/gh/SAP/cloud-mta-build-tool.svg?style=svg&circle-token=ecedd1dce3592adcd72ee4c61481972c32dcfad7)](https://circleci.com/gh/SAP/cloud-mta-build-tool)
[![Go Report Card](https://goreportcard.com/badge/github.com/SAP/cloud-mta-build-tool)](https://goreportcard.com/report/github.com/SAP/cloud-mta-build-tool)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3400/badge)](https://bestpractices.coreinfrastructure.org/projects/3400) 
[![Coverage Status](https://coveralls.io/repos/github/SAP/cloud-mta-build-tool/badge.svg?branch=cover)](https://coveralls.io/github/SAP/cloud-mta-build-tool?branch=cover)
![Beta](https://img.shields.io/badge/version-v1-green)
[![GitHub stars](https://img.shields.io/badge/contributions-welcome-orange.svg)](docs/docs/process.md)
[![dependabot](https://badgen.net/badge/Dependabot/enabled/green?icon=dependabot)](https://dependabot.com/)
[![REUSE status](https://api.reuse.software/badge/github.com/SAP/cloud-mta-build-tool)](https://api.reuse.software/info/github.com/SAP/cloud-mta-build-tool)


### Description

### Multi-Target Application

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

#### Install

For convenience the `mbt` executable is available via npmjs.com so consumers using a nodejs runtime can simply run:
- `npm install -g mbt@version`
- For possible versions see: https://www.npmjs.com/package/mbt?activeTab=versions

It is also possible to download and "install" the `mbt` executable via github releases.
- See: https://github.com/SAP/cloud-mta-build-tool/releases.

#### The Cloud MTA Build Tool Images

We supply several images for **CI environment** containing the Cloud MTA Build Tool.
The images are hosted at [Github container registry](https://github.com/orgs/SAP/packages?tab=packages&q=mbtci-)
and also at [Docker Hub registry](https://hub.docker.com/search?q=mbtci-&type=image).

The images are built from a template docker file which depends on most common technologies (Java and Node) as follows:
* [mbtci-java8-node14](https://hub.docker.com/r/devxci/mbtci-java8-node14) is built from [Dockerfile_mbtci_template](https://github.com/SAP/cloud-mta-build-tool/blob/master/Dockerfile_mbtci_template) using Node 14.
* [mbtci-java8-node16](https://hub.docker.com/r/devxci/mbtci-java8-node16) is built from [Dockerfile_mbtci_template](https://github.com/SAP/cloud-mta-build-tool/blob/master/Dockerfile_mbtci_template) using Node 16.
* [mbtci-java8-node18](https://hub.docqker.com/r/devxci/mbtci-java8-node18) is built from [Dockerfile_mbtci_template](https://github.com/SAP/cloud-mta-build-tool/blob/master/Dockerfile_mbtci_template) using Node 18.

* [mbtci-java11-node14](https://hub.docker.com/r/devxci/mbtci-java11-node14) is built from [Dockerfile_mbtci_template](https://github.com/SAP/cloud-mta-build-tool/blob/master/Dockerfile_mbtci_template) using Node 14.
* [mbtci-java11-node16](https://hub.docker.com/r/devxci/mbtci-java11-node16) is built from [Dockerfile_mbtci_template](https://github.com/SAP/cloud-mta-build-tool/blob/master/Dockerfile_mbtci_template) using Node 16.
* [mbtci-java11-node18](https://hub.docker.com/r/devxci/mbtci-java11-node18) is built from [Dockerfile_mbtci_template](https://github.com/SAP/cloud-mta-build-tool/blob/master/Dockerfile_mbtci_template) using Node 18.

And so on. 

##### How to pull the images

You should choose the relevant image type from following list to replace the `<TYPE>` template in the command/FROM according your MTA project technologies:
* java8-node14
* java8-node16
* java8-node18
* java11-node14
* java11-node16
* java11-node18
* java17-node14
* java17-node16
* java17-node18
* java19-node14
* java19-node16
* java19-node18

From the command line:
```shell
$ docker pull devxci/mbtci-<TYPE>:latest
``` 
or 
```shell
$ docker pull ghcr.io/sap/mbtci-<TYPE>:latest
```

From Dockerfile as a base image: 
```shell
FROM devxci/mbtci-<TYPE>:latest
```
or
```shell
FROM ghcr.io/sap/mbtci-<TYPE>:latest
```

E.g. if your MTA project uses Java 11 and Node 14 then you should pull the relevant image as follows: 

From the command line: 
```shell
$ docker pull devxci/mbtci-java11-node14:latest
```
or
```shell
$ docker pull ghcr.io/sap/mbtci-java11-node14:latest
```

From Dockerfile as a base image: 
```
FROM devxci/mbtci-java11-node14:latest
```
or
```
FROM ghcr.io/sap/mbtci-java11-node14:latest
```

##### How to use the images

You should choose the relevant image type from following list to replace the `<TYPE>` template in the command according your MTA project technologies:
* java8-node14
* java8-node16
* java8-node18
* java11-node14
* java11-node16
* java11-node18
* java17-node14
* java17-node16
* java17-node18
* java19-node14
* java19-node16
* java19-node18

On a Linux/Darwin machine you can run:
```shell
$ docker run -it --rm -v "$(pwd)/[proj-releative-path]:/project" devxci/mbtci-<TYPE>:latest mbt build -p=cf -t [target-folder-name]
```
This will build an mtar file for SAP Cloud Platform (Cloud Foundry). The folder containing the project needs to be mounted into the image at /project.

<b>Note:</b> The parameter `-p=cf` can be omitted as the build for cloud foundry is the default build, this is an example of the MBT build parameters, for further commands see MBT docs.

##### How to build the images

To build the images, you should run the following shell script: 

```shell
$ sh ./scripts/build_image <JAVA_VERSION> <NODE_VERSION> <MBT_VERSION>
```
E.g. to build the image for Java 11 and Node 14 you should run the following command: 
```shell
$ sh ./scripts/build_image 11.0.17 14.20.1 1.2.20
```

The Cloud MTA Build Tool published docker images on docker hub with a pre-configured set of runtime tools (nodejs/java/maven/...).

## Node.js v10/ECMAScript modules

More and more npm packages use ECMAScript modules instead of commonJS, for ECMAScript modules are the official standard format to package JavaScript code for reuse. From v1.2.25, we use axios instead of binwrap(which has moderate severity vulnerabilities) to download binary files, but axios only supports ECMAScript modules and can't work on Node.js v10. So since v1.2.25, mbt will not support Node.js v10 and lower versions.

## License

Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
Please note that Docker images can contain other software which may be licensed under different licenses. This License file is also included in the Docker image. For any usage of built Docker images please make sure to check the licenses of the artifacts contained in the images.

## Contributions

Contributions are greatly appreciated.
If you want to contribute, follow [the guidelines](docs/docs/process.md).

## Support

Please follow our [issue template](https://github.com/SAP/cloud-mta-build-tool/issues/new/choose) on how to report an issue.
