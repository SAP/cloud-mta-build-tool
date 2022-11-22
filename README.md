
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

#### Install

For convenience the `mbt` executable is available via npmjs.com so consumers using a nodejs runtime can simply run:
- `npm install -g mbt@version`
- For possible versions see: https://www.npmjs.com/package/mbt?activeTab=versions

It is also possible to download and "install" the `mbt` executable via github releases.
- See: https://github.com/SAP/cloud-mta-build-tool/releases.

#### The Cloud MTA Build Tool Images

The Cloud MTA Build Tool published docker images on docker hub with a pre-configured set of runtime tools (nodejs/java/maven/...).

##### License

Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
Please note that Docker images can contain other software which may be licensed under different licenses. This License file is also included in the Docker image. For any usage of built Docker images please make sure to check the licenses of the artifacts contained in the images.

## Contributions

Contributions are greatly appreciated.
If you want to contribute, follow [the guidelines](docs/docs/process.md).

## Support

Please follow our [issue template](https://github.com/SAP/cloud-mta-build-tool/issues/new/choose) on how to report an issue.
